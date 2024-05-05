package model

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
