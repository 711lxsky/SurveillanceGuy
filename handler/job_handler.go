package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"

	"surveillance-guy/config"
	"surveillance-guy/model"
	"surveillance-guy/util"
)

// ADDJob
// @Summary 添加定时任务
// @Description 创建一个新的定时任务并将其添加至数据库及cron调度器
// @Tags 定时任务管理
// @Accept json
// @Produce json
// @Param job body model.Job true "定时任务详情"
// @Success 200 {object} gin.H "定时任务创建成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "该任务已经存在，请勿重复添加"
// @Failure 500 {object} gin.H "任务添加失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "在调度器中创建定时任务失败" "reason" string "错误原因"
// @Router /jobs/add [post]
func ADDJob(context *gin.Context) {
	// 提取参数
	var job model.Job
	if err := context.BindJSON(&job); err != nil {
		// 解析失败
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 判断任务是否已经存在
	if util.JobIsExistInDataBaseByName(job.Name) {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage: config.JobAlreadyExistZH,
			})
		return
	}
	// 没有该任务， 写入数据库
	err := config.DataBase.Create(&job).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobAddFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	glog.Info(job)
	// 添加定时任务到 cron 调度器
	jobRun := util.JobRun{Job: job}
	jobEntryID, err := config.Cron.AddJob(job.Cron, jobRun)
	if err != nil {
		// 因为添加任务失败， 所以需要重新恢复数据， 即 revert
		err = config.DataBase.Unscoped().Delete(&job).Error
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobAddInCronFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	job.EntryID = int(jobEntryID)
	util.PrintAllJobs()
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.JobAddSuccessZH,
		})
}

// DeleteJob
// @Summary 删除定时任务
// @Description 根据提供的任务ID删除数据库中的定时任务记录及从cron调度器中移除
// @Tags 定时任务管理
// @Accept json
// @Produce json
// @Param job body model.Job true "定时任务ID"
// @Success 200 {object} gin.H "定时任务删除成功"
// @Failure 500 {object} gin.H "JSON解析错误" "reason" string "错误原因"
// @Failure 500 {object} gin.H "获取数据库中由指定 ID 指定的定时任务失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "定时任务删除失败" "reason" string "错误原因"
// @Router /jobs/delete [delete]
func DeleteJob(context *gin.Context) {
	var job model.Job
	if err := context.BindJSON(&job); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 获取指定  ID 的定时任务存储在数据库中的 EntryID, 因为请求传递过来的 EntryID 不一定正确，不充分相信用户
	glog.Info(job.ID)
	jobEntryID, err := util.GetJobEntryIDByJobID(job.ID)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobEntryIDGetFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 手动软删除
	timeNow := time.Now()
	err = config.DataBase.Model(&job).Updates(
		model.Job{
			Name: job.Name + timeNow.String(),
			Model: gorm.Model{
				DeletedAt: &timeNow},
			EntryID: 0,
			Status:  1,
		}).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobDeleteFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 在调度器中删除该任务
	config.Cron.Remove(cron.EntryID(jobEntryID))
	util.PrintAllJobs()
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.JobDeleteSuccessZH,
		})
}

// UpdateJob
// @Summary 更新定时任务
// @Description 根据提供的任务ID更新数据库中的定时任务记录及cron调度器中的任务配置
// @Tags 定时任务管理
// @Accept json
// @Produce json
// @Param job body model.Job true "定时任务详情"
// @Success 200 {object} gin.H "定时任务更新成功"
// @Failure 500 {object} gin.H "JSON解析失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "该任务不存在，请核验"
// @Failure 500 {object} gin.H "获取数据库中由指定 ID 指定的定时任务失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "在调度器中更新定时任务失败" "reason" string "错误原因"
// @Failure 500 {object} gin.H "数据库中定时任务信息更新失败" "reason" string "错误原因"
// @Router /jobs/update [put]
func UpdateJob(context *gin.Context) {
	var (
		job model.Job
		err error
	)
	if err = context.BindJSON(&job); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JSONParseErrorZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	// 判断任务是否已经存在
	if !util.JobIsExistInDataBaseByJobID(job.ID) {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage: config.JobNotExistZH,
			})
		return
	}
	// 获取指定  ID 的定时任务存储在数据库中的 EntryID, 因为请求传递过来的 EntryID 不一定正确，不充分相信用户
	glog.Info(job.ID)
	jobEntryID, err := util.GetJobEntryIDByJobID(job.ID)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobEntryIDGetFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	job.EntryID = jobEntryID
	glog.Info(job.EntryID, job.Status)
	// 在调度器中更新对应任务
	config.Cron.Remove(cron.EntryID(jobEntryID))
	if job.Status == 0 {
		jobRun := util.JobRun{Job: job}
		newJobEntryID, err := config.Cron.AddJob(job.Cron, jobRun)
		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					config.ResponseMessage:     config.JobUpdateInCronFailZH,
					config.ResponseErrorReason: err.Error(),
				})
			return
		}
		job.EntryID = int(newJobEntryID)
	} else {
		job.EntryID = 0
	}
	util.PrintAllJobs()
	// 更新数据库
	err = config.DataBase.Where(config.IDEqual, job.ID).Save(&job).Error
	if err != nil {
		config.Cron.Remove(cron.EntryID(job.EntryID))
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobUpdateInDBFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.JobUpdateSuccessZH,
		})
}

// GetAllJobs
// @Summary 获取所有定时任务列表
// @Description 查询并返回数据库中所有定时任务的信息
// @Tags 定时任务管理
// @Accept */*
// @Produce json
// @Success 200 {object} gin.H "获取定时任务列表成功" "data" []Job
// @Failure 500 {object} gin.H "定时任务列表获取失败" "reason" string "错误原因"
// @Router /jobs/list [get]
func GetAllJobs(context *gin.Context) {
	var jobs []model.Job
	err := config.DataBase.Find(&jobs).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				config.ResponseMessage:     config.JobListGetFailZH,
				config.ResponseErrorReason: err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			config.ResponseMessage: config.JobListGetSuccessZH,
			config.ResponseData:    jobs,
		})
}
