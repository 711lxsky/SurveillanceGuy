package api

import (
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

var DataBase *gorm.DB

var Cron *cron.Cron

// UserAgent 在线可以查询 https://it-tool.711lxsky.cn/user-agent-parser
var UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
var Timeout = 30

var EmailSubject = "【更新提示】 %s 有变动啦！"

var BasicAuth = false
var AuthenticateSecrets = map[string]string{
	"sur-guy": "711lxsky.",
}

var LogFilePath = "/data/surveillance-guy.INFO"
