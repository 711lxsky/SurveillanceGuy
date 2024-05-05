package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"surveillance-guy/config"
	"surveillance-guy/model"
	"surveillance-guy/util"
)

// TestEmail
// @Summary 验证电子邮件账户
// @Description 根据提供的电子邮件账户信息测试其连通性和身份验证
// @Tags 邮件测试
// @Accept json
// @Produce json
// @Param account body Account true "电子邮件账户详情"
// @Success 200 {object} gin.H "电子邮件账户身份验证成功，可以发送邮件"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "无法正确解析该 Email 账户的 SMTP 服务器主机和端口" "reason" string "错误原因"
// @Failure 500 {object} gin.H "没有该 Email 账户能够匹配的主机和端口号， 请手动输入" "reason" string "错误原因"
// @Failure 500 {object} gin.H "该 Email 账户身份认证失败， 无法使用" "reason" string "错误原因"
// @Router /emails/test [post]
func TestEmail(context *gin.Context) {
	var (
		account model.Account
		err     error
	)
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 请求没有携带邮箱密码， 去数据库去拿
	if account.Password == "" {
		var tmpAccount model.Account
		err = config.DataBase.Where(config.EmailEqual, account.Email).First(&tmpAccount).Error
		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					config.ResponseMessage:     config.EmailFindInDBFailZH,
					config.ResponseErrorReason: err.Error(),
				})
			return
		}
		account.Password = tmpAccount.Password
	}
	// 测试该账户的连通性
	err = util.EmailIsValid(account)
	if err != nil {
		switch err.Error() {
		case config.ParseEmailError:
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					config.ResponseMessage:     config.EmailAnalyzeForSMTPInfoZH,
					config.ResponseErrorReason: err.Error(),
				})
		case config.SMTPInfoNotFound:
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					config.ResponseMessage:     config.EmailNoSMTPInfoMathZH,
					config.ResponseErrorReason: err.Error(),
				})
		default:
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					config.ResponseMessage:     config.EmailAuthenticateFailZH,
					config.ResponseErrorReason: err.Error(),
				})
		}
		if account.ID != 0 {
			config.DataBase.Model(&account).Update(model.AccountStatus, model.AccountInvalid)
		}
		return
	}
	// 邮箱账户可用， 更新数据库信息
	if account.ID != 0 {
		config.DataBase.Model(&account).Update(model.AccountStatus, model.AccountValid)
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.EmailAuthenticateSuccessZH,
		})
}
