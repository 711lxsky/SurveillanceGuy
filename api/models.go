package api

import (
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	Email    string `json:"email" gorm:"not null; unique"` // 邮箱号
	Password string `json:"password" gorm:"not null"`      // 邮箱密码/授权码
	SMTPHost string `json:"host" gorm:"type:varchar(256)"` // 邮箱 SMTP 服务器地址
	SMTPPort int    `json:"port" gorm:"type:int"`          // 邮箱 SMTP 服务器端口
	Status   int    `json:"status" gorm:"type:int"`        // Email 账号状态(连通性), 是否可以发送邮件
}

type Job struct {
	gorm.Model
	Name          string `json:"name" gorm:"not null; unique"`       // 任务名称
	Cron          string `json:"cron"`                               // 定时配置
	EntryID       int    `json:"entryId" gorm:"not null"`            // cron 调度器的 job id
	Url           string `json:"url" gorm:"type:varchar(512)"`       // 监控的 目标页面URL
	OldValue      string `json:"oldValue" gorm:"type:varchar(2048)"` // 任务抓取目标的旧值
	Pattern       string `json:"pattern" gorm:"type:varchar(1024)"`  // 目标页面URL的抓取规则
	PatternStatus int    `json:"patternStatus" gorm:"type:int"`      // 抓取规则的测试状态, 0: 未测试, 1: 测试通过, 2: 测试失败 3: 测试中
	Email         string `json:"email" gorm:"not null"`              // 邮件通知接收人
	Content       string `json:"content" gorm:"type:varchar(2048)"`  // 邮件通知内容
	Status        int    `json:"status" gorm:"type:int"`             // 工作运行状态, 0: 运行中, 1: 停止
}

type Template struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null; unique"`      // 模板名称
	Corn    string `json:"cron"`                              // 定时配置
	Pattern string `json:"pattern" gorm:"type:varchar(1024)"` // 抓取规则
	Content string `json:"content" gorm:"type:varchar(2048)"` // 邮件内容
}

var SMTPHost = map[string]string{
	"qq.com":       "smtp.qq.com",
	"163.com":      "smtp.163.com",
	"126.com":      "smtp.126.com",
	"139.com":      "smtp.139.com",
	"gmail.com":    "smtp.gmail.com",
	"foxmail.com":  "smtp.foxmail.com",
	"sina.com.cn":  "smtp.sina.com.cn",
	"sohu.com":     "smtp.sohu.com",
	"yahoo.com.cn": "smtp.mail.yahoo.com.cn",
	"live.com":     "smtp.live.com",
	"263.net":      "smtp.263.net",
	"263.net.cn":   "smtp.263.net.cn",
	"x263.net":     "smtp.263.net",
	"china.com":    "smtp.china.com",
	"tom.com":      "smtp.tom.com",
	"outlook.com":  "smtp.office365.com", // Outlook和Office 365
	"hotmail.com":  "smtp.live.com",      // Hotmail
	"aol.com":      "smtp.aol.com",       // AOL
	"zoho.com":     "smtp.zoho.com",      // Zoho Mail
	"mail.com":     "smtp.mail.com",      // Mail.com
	"inbox.com":    "smtp.inbox.com",     // Inbox.com
	"gmx.com":      "smtp.gmx.com",       // GMX
	"icloud.com":   "smtp.mail.me.com",   // iCloud
}

var SMTPPort = map[string]int{
	"qq.com":       465,
	"163.com":      465,
	"126.com":      465,
	"139.com":      465,
	"gmail.com":    587,
	"foxmail.com":  465,
	"sina.com.cn":  25,
	"sohu.com":     25,
	"yahoo.com.cn": 587,
	"live.com":     587,
	"263.net":      25,
	"263.net.cn":   25,
	"x263.net":     25,
	"china.com":    25,
	"tom.com":      25,
	"outlook.com":  587, // Outlook和Office 365
	"hotmail.com":  587, // Hotmail
	"aol.com":      587, // AOL
	"zoho.com":     465, // Zoho Mail
	"mail.com":     465, // Mail.com
	"inbox.com":    465, // Inbox.com
	"gmx.com":      587, // GMX
	"icloud.com":   587, // iCloud
}

type JobRun struct {
	Job Job
}

func (jobRun JobRun) Run() {
	// Run 执行定时任务
	infoPrefix := "[Job#%d][%s]"
	// 任务执行前后打印提示信息
	glog.Infof("======= [Job#%d][%s][%s][Status: %d][EntryID: %d][OldValue: %s] Start...",
		jobRun.Job.ID, jobRun.Job.Name, jobRun.Job.Cron, jobRun.Job.Status, jobRun.Job.EntryID, jobRun.Job.OldValue)
	defer glog.Infof("-------- [Job#%d][%s][%s][Status: %d][EntryID: %d][OldValue: %s] End --------",
		jobRun.Job.ID, jobRun.Job.Name, jobRun.Job.Cron, jobRun.Job.Status, jobRun.Job.EntryID, jobRun.Job.OldValue)
	// 执行定时任务
	err := WatchJob(jobRun.Job)
	if err != nil {
		glog.Errorf(infoPrefix+err.Error(), jobRun.Job.ID, jobRun.Job.Name)
	}
}
