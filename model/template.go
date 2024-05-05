package model

import (
	"github.com/jinzhu/gorm"
)

type Template struct {
	gorm.Model
	Name    string `json:"name" gorm:"not null; unique"`      // 模板名称
	Corn    string `json:"cron"`                              // 定时配置
	Pattern string `json:"pattern" gorm:"type:varchar(1024)"` // 抓取规则
	Content string `json:"content" gorm:"type:varchar(2048)"` // 邮件内容
}
