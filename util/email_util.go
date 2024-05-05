package util

import (
	"crypto/tls"
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"

	"surveillance-guy/config"
	"surveillance-guy/model"
)

// SendEmail
// 发送邮件
func SendEmail(account model.Account, maiTo []string, subject, body string) error {
	var (
		err      error
		smtpHost string
		smtpPort int
	)
	if account.SMTPHost == "" || account.SMTPPort == 0 {
		smtpHost, smtpPort, err = ParseSMTPInfoByEmail(account.Email)
	} else {
		smtpHost = account.SMTPHost
		smtpPort = account.SMTPPort
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
		return "", 0, fmt.Errorf(config.ParseEmailError)
	}
	suffix := splitRes[len(splitRes)-1]
	// 根据后缀获取SMTP 信息
	smtpHost, ok1 := model.SMTPHost[suffix]
	smtpPort, ok2 := model.SMTPPort[suffix]
	if ok1 && ok2 {
		return smtpHost, smtpPort, nil
	} else {
		return smtpHost, smtpPort, fmt.Errorf(config.SMTPInfoNotFound)
	}
}

// EmailIsValid
// 测试邮箱的有效性， 是否可以连通
func EmailIsValid(account model.Account) error {
	var (
		smtpHost string
		smtpPort int
		err      error
	)
	if account.SMTPHost == "" || account.SMTPPort == 0 {
		// 未提供SMTP 信息, 根据邮箱后缀解析
		smtpHost, smtpPort, err = ParseSMTPInfoByEmail(account.Email)
	} else {
		// 提供了 SMTP 信息
		smtpHost = account.SMTPHost
		smtpPort = account.SMTPPort
	}
	if err != nil {
		return err
	}
	// 拨号， 向 SMTP 服务器进行身份验证
	dialer := gomail.NewDialer(smtpHost, smtpPort, account.Email, account.Password)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	_, err = dialer.Dial()
	if err != nil {
		return err
	}
	return nil
}
