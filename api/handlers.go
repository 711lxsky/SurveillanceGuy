package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// AuthenticateHandler
// 认证检测
func AuthenticateHandler(context *gin.Context) {
	userName := context.MustGet(gin.AuthUserKey).(string)
	if _, ok := AuthenticateSecrets[userName]; ok {
		context.JSON(http.StatusOK, gin.H{"message": "认证成功", gin.AuthUserKey: userName})
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "认证失败, 该账号不存在", gin.AuthUserKey: userName})
	}
}
