package model

import (
	"github.com/jinzhu/gorm"
)

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

var (
	ID            = "id"
	URL           = "url"
	Pattern       = "pattern"
	Type          = "type"
	RE            = "re"
	PatternStatus = "patten_status"
)

var (
	RegexPatternValid = 1
)
