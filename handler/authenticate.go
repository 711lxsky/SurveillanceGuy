package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"surveillance-guy/config"
)

// AuthenticateHandler
// @Summary 用户认证接口
// @Description 验证用户凭据并返回认证结果
// @Tags 用户认证
// @Accept json
// @Produce json
// @Param username path string true "用户名" description "用户的唯一标识"
// @Success 200 {object} gin.H "认证成功"
// @Failure 401 {object} gin.H "认证失败, 该账号不存在"
// @Router /authenticate [post]
func AuthenticateHandler(context *gin.Context) {
	userName := context.MustGet(gin.AuthUserKey).(string)
	if _, ok := config.AuthenticateSecrets[userName]; ok {
		context.JSON(http.StatusOK, gin.H{"message": config.AuthenticateSuccessZH, gin.AuthUserKey: userName})
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"message": config.AuthenticateFailZH, gin.AuthUserKey: userName})
	}
}
