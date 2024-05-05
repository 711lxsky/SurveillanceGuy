package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"surveillance-guy/config"
	"surveillance-guy/model"
)

// AddTemplate
// @Summary 创建任务模板
// @Description 添加一个新的任务模板到数据库
// @Tags 任务模板管理
// @Accept json
// @Produce json
// @Param template body model.Template true "任务模板详情"
// @Success 200 {object} gin.H "任务模板创建成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "任务模板创建失败" "reason" string "错误原因"
// @Router /templates/add [post]
func AddTemplate(context *gin.Context) {
	var template model.Template
	if err := context.BindJSON(&template); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	err := config.DataBase.Create(&template).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.TemplateAddFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.TemplateAddSuccessZH,
		})
}

// DeleteTemplate
// @Summary 删除任务模板
// @Description 根据提供的任务模板信息从数据库中软删除
// @Tags 任务模板管理
// @Accept json
// @Produce json
// @Param template body model.Template true "任务模板ID"
// @Success 200 {object} gin.H "任务模板删除成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "任务模板删除失败" "reason" string "错误原因"
// @Router /templates/delete [delete]
func DeleteTemplate(context *gin.Context) {
	var template model.Template
	if err := context.BindJSON(&template); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 软删除
	timeNow := time.Now()
	err := config.DataBase.Model(&template).Updates(
		model.Template{
			Name: template.Name + timeNow.String(),
			Model: gorm.Model{
				DeletedAt: &timeNow},
		}).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.TemplateDeleteFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.TemplateDeleteSuccessZH,
		})
}

// UpdateTemplate
// @Summary 更新任务模板
// @Description 根据提供的任务模板ID更新其信息
// @Tags 任务模板管理
// @Accept json
// @Produce json
// @Param template body model.Template true "任务模板详情"
// @Success 200 {object} gin.H "任务模板更新成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "任务模板更新失败" "reason" string "错误原因"
// @Router /templates/update [put]
func UpdateTemplate(context *gin.Context) {
	var (
		template model.Template
		err      error
	)
	if err = context.BindJSON(&template); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	err = config.DataBase.Where(config.IDEqual, template.ID).Save(&template).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.TemplateUpdateFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.TemplateUpdateSuccessZH,
		})
}

// GetAllTemplates
// @Summary 获取所有任务模板
// @Description 查询并返回数据库中所有任务模板的信息
// @Tags 任务模板管理
// @Accept */*
// @Produce json
// @Success 200 {object} gin.H "获取任务模板成功" "data" []Template
// @Failure 500 {object} gin.H "获取任务模板失败" "reason" string "错误原因"
// @Router /templates/all [get]
func GetAllTemplates(context *gin.Context) {
	var templates []model.Template
	err := config.DataBase.Find(&templates).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.TemplateListGetFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.TemplateListGetSuccessZH,
			config.ResponseData:    templates,
		})
}
