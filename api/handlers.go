package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
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

// LogTail
// 持续输出日志内容
func LogTail(context *gin.Context) {
	// 升级 GET 请求为 Websocket 协议
	websocket, err := upGrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer websocket.Close()
	// 监控日志文件
	tails, err := tail.TailFile(LogFilePath, tail.Config{
		ReOpen:    true,
		Follow:    true,
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "日志文件打开失败",
				"reason":  err.Error(),
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
				fmt.Printf("tail fole close reopen, filename: %s\n", tails.Filename)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			fmt.Println("msg:", msg)
			// 写入 Websocket 数据
			err = websocket.WriteMessage(messageType, []byte(msg.Text))
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}

// ADDJob
// 添加定时任务
func ADDJob(context *gin.Context) {
	// 提取参数
	var job Job
	if err := context.BindJSON(&job); err != nil {
		// 解析失败
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	// 判断任务是否已经存在
	if JobIsExistInDataBaseByName(job.Name) {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "该任务已经存在， 请勿重复添加",
				"reason":  "",
			})
		return
	}
	// 没有该任务， 写入数据库
	err := DataBase.Create(&job).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "任务添加失败",
				"reason":  err.Error(),
			})
		return
	}
	glog.Info(job)
	// 添加定时任务到 cron 调度器
	jobRun := JobRun{Job: job}
	jobEntryID, err := Cron.AddJob(job.Cron, jobRun)
	if err != nil {
		// 因为添加任务失败， 所以需要重新恢复数据， 即 revert
		err = DataBase.Unscoped().Delete(&job).Error
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "在调度器中创建定时任务失败",
				"reason":  err.Error(),
			})
		return
	}
	job.EntryID = int(jobEntryID)
	PrintAllJobs()
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "定时任务创建成功",
		})
}

func DeleteJob(context *gin.Context) {
	var job Job
	if err := context.BindJSON(&job); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析错误",
				"reason":  err.Error(),
			})
		return
	}
	// 获取指定  ID 的定时任务存储在数据库中的 EntryID, 因为请求传递过来的 EntryID 不一定正确，不充分相信用户
	glog.Info(job.ID)
	jobEntryID, err := GetJobEntryIDByJobID(job.ID)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "获取数据库中由指定 ID 指定的定时任务失败",
				"reason":  err.Error(),
			})
		return
	}
	// 手动软删除
	timeNow := time.Now()
	err = DataBase.Model(&job).Updates(
		Job{
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
				"message": "定时任务删除失败",
				"reason":  err.Error(),
			})
		return
	}
	// 在调度器中删除该任务
	Cron.Remove(cron.EntryID(jobEntryID))
	PrintAllJobs()
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "定时任务删除成功",
		})
}

// UpdateJob
// 更新定时任务
func UpdateJob(context *gin.Context) {
	var (
		job Job
		err error
	)
	if err = context.BindJSON(&job); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	// 判断任务是否已经存在
	if !JobIsExistInDataBaseByJobID(job.ID) {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "该任务不存在， 请核验",
				"reason":  "",
			})
		return
	}
	// 获取指定  ID 的定时任务存储在数据库中的 EntryID, 因为请求传递过来的 EntryID 不一定正确，不充分相信用户
	glog.Info(job.ID)
	jobEntryID, err := GetJobEntryIDByJobID(job.ID)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "获取数据库中由指定 ID 指定的定时任务失败",
				"reason":  err.Error(),
			})
		return
	}
	job.EntryID = jobEntryID
	glog.Info(job.EntryID, job.Status)
	// 在调度器中更新对应任务
	Cron.Remove(cron.EntryID(jobEntryID))
	if job.Status == 0 {
		jobRun := JobRun{Job: job}
		newJobEntryID, err := Cron.AddJob(job.Cron, jobRun)
		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "在调度器中更新定时任务失败",
					"reason":  err.Error(),
				})
			return
		}
		job.EntryID = int(newJobEntryID)
	} else {
		job.EntryID = 0
	}
	PrintAllJobs()
	// 更新数据库
	err = DataBase.Where("id = ?", job.ID).Save(&job).Error
	if err != nil {
		Cron.Remove(cron.EntryID(job.EntryID))
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "数据库中定时任务信息更新失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "定时任务更新成功",
		})
}

// GetAllJobs
// 获取数据库中所有定时任务
func GetAllJobs(context *gin.Context) {
	var jobs []Job
	err := DataBase.Find(&jobs).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "定时任务列表获取失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "获取定时任务列表成功",
			"data":    jobs,
		})
}

// AddAccount
// 新建邮件通知账户
func AddAccount(context *gin.Context) {
	var account Account
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	err := DataBase.Create(&account).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "邮件通知账户新建失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "邮件通知账户创建成功",
		})
}

// DeleteAccount
// 删除邮件通知账户
func DeleteAccount(context *gin.Context) {
	var account Account
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	// 软删除
	timeNow := time.Now()
	err := DataBase.Model(&account).Updates(
		Account{
			Email: account.Email + timeNow.String(),
			Model: gorm.Model{
				DeletedAt: &timeNow},
		}).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "邮件通知账户删除失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "邮件通知账户删除成功",
		})
}

// UpdateAccount
// 更新邮件通知账户
func UpdateAccount(context *gin.Context) {
	var account Account
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	// 根据账户 ID 拿到信息并更新
	err := DataBase.Where("id = ?", account.ID).Save(&account).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "邮件通知账户更新失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "邮件通知账户更新成功",
		})
}

// GetAllAccounts
// 获取所有邮件通知账户
func GetAllAccounts(context *gin.Context) {
	var accounts []Account
	err := DataBase.Find(&accounts).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "邮件通知账户列表获取失败",
				"reason":  err.Error(),
			})
		return
	}
	// 擦除 password 字段
	n := len(accounts)
	for i := 0; i < n; i++ {
		accounts[i].Password = "********"
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "获取邮件通知账户列表成功",
			"data":    accounts,
		})
}

// TestRegexPattern
// 测试正则表达式效果
func TestRegexPattern(context *gin.Context) {
	var (
		testRes string
		job     Job
		err     error
	)
	jobID := context.Query("id")
	url := context.Query("url")
	pattern := context.Query("pattern")
	patternType := context.DefaultQuery("type", "re")
	DataBase.Where("id = ?", jobID).First(&job)
	html, err := GetHtmlByUrl(url)
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "获取页面 html 源码失败",
				"reason":  err.Error(),
			})
		return
	}
	// 根据 pattern 对 html 源码进行匹配
	switch patternType {
	case "re":
		testRes, err = MatchTargetByRegexPattern(html, pattern)
	default:
		err = fmt.Errorf("Pattern-type `%s` is not found", patternType)
	}
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "抓取规则无效",
				"reason":  err.Error(),
			})
		return
	}
	// 抓取规则有效， 更新数据库， 将状态设置为有效, 即测试通过
	if job.ID != 0 {
		DataBase.Model(&job).Update("patten_status", 1)
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "正则表达式测试成功， 匹配到内容",
			"data":    testRes,
		})
}

// TestEmail
// 测试邮箱是否可以发送邮箱
func TestEmail(context *gin.Context) {
	var (
		account Account
		err     error
	)
	if err := context.BindJSON(&account); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	// 请求没有携带邮箱密码， 去数据库去拿
	if account.Password == "" {
		var tmpAccount Account
		err = DataBase.Where("email = ?", account.Email).First(&tmpAccount).Error
		if err != nil {
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "数据库中未找到该邮箱",
					"reason":  err.Error(),
				})
			return
		}
		account.Password = tmpAccount.Password
	}
	// 测试该账户的连通性
	err = EmailIsValid(account)
	if err != nil {
		switch err.Error() {
		case ParseEmailError:
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "无法正确解析该 Email 账户的 SMTP 服务器主机和端口",
					"reason":  err.Error(),
				})
		case SMTPInfoNotFound:
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "没有该 Email 账户能够匹配的主机和端口号， 请手动输入",
					"reason":  err.Error(),
				})
		default:
			context.AbortWithStatusJSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "该 Email 账户身份认证失败， 无法使用",
					"reason":  err.Error(),
				})
		}
		if account.ID != 0 {
			DataBase.Model(&account).Update("status", 2)
		}
		return
	}
	// 邮箱账户可用， 更新数据库信息
	if account.ID != 0 {
		DataBase.Model(&account).Update("status", 1)
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "该 Email 账户身份验证成功， 可以发送邮件",
		})
}

// AddTemplate
// 添加任务模板
func AddTemplate(context *gin.Context) {
	var template Template
	if err := context.BindJSON(&template); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	err := DataBase.Create(&template).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "任务模板创建失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "任务模板创建成功",
		})
}

// DeleteTemplate
// 删除任务模板
func DeleteTemplate(context *gin.Context) {
	var template Template
	if err := context.BindJSON(&template); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	// 软删除
	timeNow := time.Now()
	err := DataBase.Model(&template).Updates(
		Template{
			Name: template.Name + timeNow.String(),
			Model: gorm.Model{
				DeletedAt: &timeNow},
		}).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "任务模板删除失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "任务模板删除成功",
		})
}

// UpdateTemplate
// 更新任务模板
func UpdateTemplate(context *gin.Context) {
	var (
		template Template
		err      error
	)
	if err = context.BindJSON(&template); err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "JSON 解析失败",
				"reason":  err.Error(),
			})
		return
	}
	err = DataBase.Where("id = ?", template.ID).Save(&template).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "任务模板更新失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "任务模板更新成功",
		})
}

// GetAllTemplates
// 获取所有任务模板
func GetAllTemplates(context *gin.Context) {
	var templates []Template
	err := DataBase.Find(&templates).Error
	if err != nil {
		context.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "获取任务模板失败",
				"reason":  err.Error(),
			})
		return
	}
	context.JSON(
		http.StatusOK,
		gin.H{
			"message": "获取任务模板成功",
			"data":    templates,
		})
}
