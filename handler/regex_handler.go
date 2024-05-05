package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"surveillance-guy/config"
	"surveillance-guy/model"
	"surveillance-guy/util"
)

// TestRegexPattern
// @Summary 测试正则表达式匹配
// @Description 根据提供的URL、正则表达式和类型对页面HTML源码进行匹配测试
// @Tags 正则测试
// @Accept */*
// @Produce json
// @Param id query string true "任务ID"
// @Param url query string true "页面URL"
// @Param pattern query string true "正则表达式"
// @Param type query string false "匹配类型，默认为're'（正则）" default(re)
// @Success 200 {object} gin.H "正则表达式测试成功，匹配到内容" "data" string
// @Failure 500 {object} gin.H "获取页面 html 源码失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "抓取规则无效" "reason" string "错误原因"
// @Router /test-email [get]测试正则表达式效果
func TestRegexPattern(context *gin.Context) {
	var (
		testRes string
		job     model.Job
		err     error
	)
	jobID := context.Query(model.ID)
	url := context.Query(model.URL)
	pattern := context.Query(model.Pattern)
	patternType := context.DefaultQuery(model.Type, model.RE)
	config.DataBase.Where(config.IDEqual, jobID).First(&job)
	html, err := util.GetHtmlByUrl(url)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.HtmlCodeGetFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 根据 pattern 对 html 源码进行匹配
	switch patternType {
	case model.RE:
		testRes, err = util.MatchTargetByRegexPattern(html, pattern)
	default:
		err = fmt.Errorf(config.PatternTypeNotFound, patternType)
	}
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.RegexPatternInvalidZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 抓取规则有效， 更新数据库， 将状态设置为有效, 即测试通过
	if job.ID != 0 {
		config.DataBase.Model(&job).Update(model.PatternStatus, model.RegexPatternValid)
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.RegexPatternValidZH,
			config.ResponseData:    testRes,
		})
}
