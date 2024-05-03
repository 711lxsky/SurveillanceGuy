package api

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
	"gopkg.in/gomail.v2"
)

func WatchJob(job Job) error {
	// 爬取目标页面指定内容, 和数据库中对比, 如果有变动, 发送邮件通知
	// 定时任务结束时输出缓冲区日志
	defer glog.Flush()

	var jobOldValue, jobNewValue string
	var infoPrefix = "[Job#%d][%s]"

	// 爬取目标页面 html
	glog.Infof(infoPrefix+"Crawling the target page...", job.ID, job.Name)
	html, err := GetHtmlByUrl(job.Url)
	if err != nil {
		return err
	}
	// 匹配指定内容, 获取新值
	glog.Infof(infoPrefix+"Matching the specified content...", job.ID, job.Name)
	// 根据正则表达式拿到对应内容
	jobNewValue, err = MatchTargetByRegexPattern(html, job.Pattern)
	if err != nil {
		return err
	}
	glog.Infof(infoPrefix+"Got the new value: %s", job.ID, job.Name, jobNewValue)
	// 从数据库取出旧值
	glog.Infof(infoPrefix+"Getting the old value from Database...", job.ID, job.Name)
	tmpJob := Job{}
	err = DataBase.First(&tmpJob, job.ID).Error
	if err != nil {
		return err
	}
	jobOldValue = tmpJob.OldValue
	glog.Infof(infoPrefix+"Got the old value: %s", job.ID, job.Name, jobOldValue)
	// 判断新旧值是否相同
	glog.Infof(infoPrefix+"Comparing the new value: '%s' and the old value: '%s'...", job.ID, job.Name, jobNewValue, jobOldValue)
	if jobNewValue == jobOldValue {
		// 相同, 不管
		glog.Infof(infoPrefix+"The new value is the same as the old value, no need to send email, skipping", job.ID, job.Name)
	} else {
		// 不同, 更新数据库, 发送通知
		glog.Infof(infoPrefix+"The new value is different from the old value, updating and sending email...", job.ID, job.Name)
		// 更形
		err = DataBase.Model(&job).Update("old_value", jobNewValue).Error
		if err != nil {
			return err
		}
		var account Account
		err = DataBase.Where("email = ?", job.Email).First(&account).Error
		if err != nil {
			return err
		}
		// 替换 job.Content 内容中的变量
		compileTargetRes, _ := regexp.Compile("%target%")
		job.Content = compileTargetRes.ReplaceAllLiteralString(job.Content, jobNewValue)
		compileNameRes, _ := regexp.Compile("%name%")
		job.Content = compileNameRes.ReplaceAllLiteralString(job.Content, job.Name)
		// 发送通知
		err = SendEmail(account, []string{job.Email}, fmt.Sprintf(EmailSubject, job.Name), job.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetHtmlByUrl
// 抓取指定 url 的 html 页面源码
func GetHtmlByUrl(url string) ([]byte, error) {
	// 生成 Client 客户端
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: time.Duration(Timeout) * time.Second,
	}
	// 构建请求
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}
	request.Header.Add("User-Agent", UserAgent)
	// 发起请求
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}
	// 关闭响应体
	defer response.Body.Close()
	html, err := DataEncoding(response.Body)
	if err != nil {
		return []byte{}, err
	}
	return html, nil
}

// DataEncoding
// 自动转换页面编码, html 页面本身决定编码
func DataEncoding(r io.Reader) ([]byte, error) {
	oldReader := bufio.NewReader(r)
	bytes, err := oldReader.Peek(1024)
	if err != nil {
		return []byte{}, err
	}
	encoding, _, _ := charset.DetermineEncoding(bytes, "")
	reader := transform.NewReader(oldReader, encoding.NewDecoder())

	// 读取相应体
	html, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	return html, nil
}

// MatchTargetByRegexPattern
// 根据正则表达式匹配目标内容
func MatchTargetByRegexPattern(content []byte, pattern string) (string, error) {
	compileRes, _ := regexp.Compile(pattern)
	items := compileRes.FindSubmatch(content)
	if len(items) >= 2 {
		res := string(items[1])
		return res, nil
	} else {
		return "", fmt.Errorf("Can't match target")
	}
}

// SendEmail
// 发送邮件
func SendEmail(account Account, maiTo []string, subject, body string) error {
	var (
		err      error
		smtpHost string
		smtpPort int
	)
	if account.STMPHost == "" || account.STMPPort == 0 {
		smtpHost, smtpPort, err = ParseSMTPInfoByEmail(account.Email)
	} else {
		smtpHost = account.STMPHost
		smtpPort = account.STMPPort
	}
	if err != nil {
		return err
	}
	// 构建邮件
	newEmailMessage := gomail.NewMessage()
	newEmailMessage.SetHeader("Form", account.Email)
	newEmailMessage.SetHeader("To", maiTo...)
	newEmailMessage.SetHeader("Subject", subject)
	newEmailMessage.SetBody("text/html", body)
	// 发送邮件
	dialer := gomail.NewDialer(smtpHost, smtpPort, account.Email, account.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err = dialer.DialAndSend(newEmailMessage)
	return err
}

// ParseSMTPInfoByEmail
// 根据邮箱后缀解析 SMTP 信息
func ParseSMTPInfoByEmail(email string) (string, int, error) {
	splitRes := strings.Split(email, "@")
	if len(splitRes) <= 1 {
		return "", 0, fmt.Errorf("Can't parse email suffix")
	}
	suffix := splitRes[len(splitRes)-1]
	// 根据后缀获取SMTP 信息
	smtpHost, ok1 := SMTPHost[suffix]
	smtpPort, ok2 := SMTPPort[suffix]
	if ok1 && ok2 {
		return smtpHost, smtpPort, nil
	} else {
		return smtpHost, smtpPort, fmt.Errorf("Can't found target SMTP information for the email-suffix")
	}
}
