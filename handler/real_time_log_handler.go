package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"

	"surveillance-guy/config"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// LogTail
// @Summary 日志实时推送接口
// @Description 通过WebSocket协议获取实时日志流
// @Tags 日志监控
// @Accept */*
// @Produce */*
// @Router /websocket [get]
func LogTail(context *gin.Context) {
	// 升级 GET 请求为 Websocket 协议
	websocket, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer websocket.Close()
	// 监控日志文件
	tails, err := tail.TailFile(config.LogFilePath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.LogFileOpenFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	for {
		// 读取 ws 中数据， 接收客户端连接
		messageType, messageContent, err := websocket.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(string(messageContent))
		// 开始
		var (
			msg *tail.Line
			ok  bool
		)
		for {
			msg, ok = <-tails.Lines
			if !ok {
				fmt.Printf(config.TailFileReopen, tails.Filename)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			fmt.Println(config.ResponseMessage, msg)
			// 写入 Websocket 数据
			err = websocket.WriteMessage(messageType, []byte(msg.Text))
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}
