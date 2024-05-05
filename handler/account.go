package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"surveillance-guy/config"
	"surveillance-guy/model"
)

// AddAccount
// @Summary 新建邮件通知账户
// @Description 创建一个新的邮件通知账户并将其添加至数据库
// @Tags 邮件通知管理
// @Accept json
// @Produce json
// @Param account body Account true "邮件通知账户详情"
// @Success 200 {object} gin.H "邮件通知账户创建成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "邮件通知账户新建失败" "reason" string "错误原因"
// @Router /accounts/add [post]
func AddAccount(context *gin.Context) {
	var account model.Account
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	err := config.DataBase.Create(&account).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.AccountAddFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.AccountAddSuccessZH,
		})
}

// DeleteAccount
// @Summary 删除邮件通知账户
// @Description 根据提供的账户信息从数据库中软删除邮件通知账户
// @Tags 邮件通知管理
// @Accept json
// @Produce json
// @Param account body Account true "邮件通知账户ID"
// @Success 200 {object} gin.H "邮件通知账户删除成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "邮件通知账户删除失败" "reason" string "错误原因"
// @Router /accounts/delete [delete]
func DeleteAccount(context *gin.Context) {
	var account model.Account
	if err := context.BindJSON(&account); err != nil {
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
	err := config.DataBase.Model(&account).Updates(
		model.Account{
			Email: account.Email + timeNow.String(),
			Model: gorm.Model{
				DeletedAt: &timeNow},
		}).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.AccountDeleteFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.AccountDeleteSuccessZH,
		})
}

// UpdateAccount
// @Summary 更新邮件通知账户
// @Description 根据提供的账户ID更新数据库中的邮件通知账户信息
// @Tags 邮件通知管理
// @Accept json
// @Produce json
// @Param account body Account true "邮件通知账户详情"
// @Success 200 {object} gin.H "邮件通知账户更新成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "邮件通知账户更新失败" "reason" string "错误原因"
// @Router /accounts/update [put]更新邮件通知账户
func UpdateAccount(context *gin.Context) {
	var account model.Account
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 根据账户 ID 拿到信息并更新
	err := config.DataBase.Where(config.IDEqual, account.ID).Save(&account).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.AccountUpdateFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.AccountUpdateSuccessZH,
		})
}

// GetAllAccounts
// @Summary 获取所有邮件通知账户列表
// @Description 查询并返回数据库中所有邮件通知账户的信息，敏感信息如密码会被隐藏
// @Tags 邮件通知管理
// @Accept */*
// @Produce json
// @Success 200 {object} gin.H "获取邮件通知账户列表成功" "data" []Account
// @Failure 500 {object} gin.H "邮件通知账户列表获取失败" "reason" string "错误原因"
// @Router /accounts/list [get]
func GetAllAccounts(context *gin.Context) {
	var accounts []model.Account
	err := config.DataBase.Find(&accounts).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.AccountListGetFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 擦除 password 字段
	n := len(accounts)
	for i := 0; i < n; i++ {
		accounts[i].Password = config.PasswordEncoded
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.AccountListGetSuccessZH,
			config.ResponseData:    accounts,
		})
}
