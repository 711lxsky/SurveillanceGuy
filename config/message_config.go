package config

var (
	ParseEmailError     = "Can't parse email suffix"
	SMTPInfoNotFound    = "Can't found target SMTP information for the email-suffix"
	TailFileReopen      = "tail file close reopen, filename: %s\n"
	PatternTypeNotFound = "Pattern-type `%s` is not found"
)

var (
	JSONParseErrorZH           = " JSON 解析失败"
	AccountAddFailZH           = "邮箱通知账户添加失败"
	AccountAddSuccessZH        = "邮箱通知账户添加成功"
	AccountDeleteFailZH        = "邮箱通知账户删除失败"
	AccountDeleteSuccessZH     = "邮箱通知账户删除成功"
	AccountUpdateFailZH        = "邮箱通知账户更新失败"
	AccountUpdateSuccessZH     = "邮箱通知账户更新成功"
	AccountListGetFailZH       = "邮箱通知账户列表获取失败"
	AccountListGetSuccessZH    = "邮箱通知账户列表获取成功"
	AuthenticateSuccessZH      = "认证成功"
	AuthenticateFailZH         = "认证失败, 该账号不存在"
	EmailFindInDBFailZH        = "数据库中未找到该邮箱"
	EmailAnalyzeForSMTPInfoZH  = "无法正确解析该 Email 账户的 SMTP 服务器主机和端口"
	EmailNoSMTPInfoMathZH      = "没有该 Email 账户能够匹配的主机和端口号， 请手动输入"
	EmailAuthenticateFailZH    = "该 Email 账户身份认证失败， 无法使用"
	EmailAuthenticateSuccessZH = "该 Email 账户身份验证成功， 可以发送邮件"
	JobAlreadyExistZH          = "该任务已经存在， 请勿重复添加"
	JobAddFailZH               = "任务添加失败"
	JobAddInCronFailZH         = "在调度器中创建定时任务失败"
	JobAddSuccessZH            = "定时任务添加成功"
	JobEntryIDGetFailZH        = "获取数据库中由指定 ID 指定的定时任务失败"
	JobDeleteFailZH            = "定时任务删除失败"
	JobDeleteSuccessZH         = "定时任务删除成功"
	JobNotExistZH              = "该定时任务不存在， 请核验"
	JobUpdateInCronFailZH      = "在调度器中更新定时任务失败"
	JobUpdateInDBFailZH        = "在数据库中更新定时任务信息失败"
	JobUpdateSuccessZH         = "定时任务更新成功"
	JobListGetFailZH           = "定时任务列表获取失败"
	JobListGetSuccessZH        = "定时任务列表获取成功"
	LogFileOpenFailZH          = "日志文件打开失败"
	HtmlCodeGetFailZH          = "获取页面 html 源码失败"
	RegexPatternInvalidZH      = "抓取规则无效"
	RegexPatternValidZH        = "正则表达式测试成功， 匹配到内容"
	TemplateAddFailZH          = "任务模板创建失败"
	TemplateAddSuccessZH       = "任务模板创建成功"
	TemplateDeleteFailZH       = "任务模板删除失败"
	TemplateDeleteSuccessZH    = "任务模板删除成功"
	TemplateUpdateFailZH       = "任务模板更新失败"
	TemplateUpdateSuccessZH    = "任务模板更新成功"
	TemplateListGetFailZH      = "获取任务模板列表失败"
	TemplateListGetSuccessZH   = "获取任务模板列表成功"
)