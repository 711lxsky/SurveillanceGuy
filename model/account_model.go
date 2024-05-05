package model

import (
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

var (
	AccountStatus  = "status"
	AccountInvalid = 2
	AccountValid   = 1
)
