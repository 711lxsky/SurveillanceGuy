package main

import (
	"flag"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/robfig/cron/v3"

	"surveillance-guy/config"
	"surveillance-guy/handler"
	"surveillance-guy/model"
	"surveillance-guy/util"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "surveillance-guy/cmd/docs"
)

// SyncJobsInDataBase
// 同步数据库中的定时任务
func SyncJobsInDataBase() error {
	jobs := []model.Job{}
	err := config.DataBase.Find(&jobs).Error
	if err != nil {
		return err
	}
	for _, job := range jobs {
		glog.Info(job)
		if job.Status != 0 || job.DeletedAt != nil {
			// 筛选一下
			continue
		}
		// 在任务调度器中创建新任务
		jobRun := model.JobRun{Job: job}
		jobNewEntryID, err := config.Cron.AddJob(job.Cron, jobRun)
		if err != nil {
			return err
		}
		err = config.DataBase.Model(&job).Update("EntryID", jobNewEntryID).Error
		if err != nil {
			return err
		}
	}
	util.PrintAllJobs()
	return nil
}

func main() {
	// 初始化日志库
	flag.Parse()
	// Flush 守护进程间隔 30s 周期性刷新缓冲区中的日志
	defer glog.Flush()
	// 连接 sqlite3 数据库
	var err error
	config.DataBase, err = gorm.Open("sqlite3", "surveillance_guy.db")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer config.DataBase.Close()
	// 自动迁移模式， 保持更新到最新
	// 仅创建表， 缺少列和索引， 不会改变现有列的类型或删除未使用的列以保护数据
	config.DataBase.AutoMigrate(&model.Account{}, &model.Job{}, &model.Template{})
	// 创建并开始 cron 调度定时任务
	config.Cron = cron.New()
	// 同步数据库中存在的定时任务
	err = SyncJobsInDataBase()
	if err != nil {
		glog.Error(err.Error())
	}
	config.Cron.Start()
	// 创建 gin 实例
	engine := gin.Default()
	// 添加 CORS 中间件， 允许跨域请求访问
	engine.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 路由绑定
	var v1 = engine.Group("/handler/v1")
	if config.BasicAuth {
		v1.Handlers = append(v1.Handlers, gin.BasicAuth(gin.Accounts(config.AuthenticateSecrets)))
	}
	{
		// 认证校验
		v1.GET("/secrets", handler.AuthenticateHandler)
		// Websocket 日志持续输出
		v1.GET("/websocket", handler.LogTail)
		// 定时任务 CRUD
		v1.POST("/job", handler.ADDJob)
		v1.DELETE("/job", handler.DeleteJob)
		v1.PUT("/job", handler.UpdateJob)
		v1.GET("/job", handler.GetAllAccounts)
		// 邮箱账号 CRUD
		v1.POST("/account", handler.AddAccount)
		v1.DELETE("/account", handler.DeleteAccount)
		v1.PUT("/account", handler.UpdateAccount)
		v1.GET("/account", handler.GetAllAccounts)
		// 功能测试接口
		v1.GET("/testpattern", handler.TestRegexPattern)
		v1.POST("/testemail", handler.TestEmail)
		// 任务模板 CRUD
		v1.POST("/template", handler.AddTemplate)
		v1.DELETE("/template", handler.DeleteTemplate)
		v1.PUT("/template", handler.UpdateTemplate)
		v1.GET("/template", handler.GetAllTemplates)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := "8848"
	engine.Run(":" + port)
}
