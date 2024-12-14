package service

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/infrastructure"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/utils"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// SendWeeklyEmail 后台多协程发送周报邮件
func SendWeeklyEmail(AdminID uint) {
	go func() {
		var wg sync.WaitGroup
		var routineNum = config.Configs.Smtp.WeeklyEmail.RoutineNum
		var userEmails []string
		var TotalSuccessCount int
		var TotalFailCount int

		// 获取超时时间
		TimeOut, err := time.ParseDuration(config.Configs.Smtp.WeeklyEmail.TimeOut)
		if err != nil {
			logger.Log.Error(nil, err.Error())
			return
		}

		// 查询所有订阅了周报的用户的 email
		if err = infrastructure.PostgresDb.Model(&model.User{}).Where("subscribe_weekly_email = true").Pluck("email", &userEmails).Error; err != nil {
			logger.Log.Error(nil, err.Error())
			return
		}

		emailsChan := make(chan string, len(userEmails))
		// 把所有的 email 放入 emailsChan里
		for _, email := range userEmails {
			emailsChan <- email
		}
		defer close(emailsChan)
		// 这两个切片用来记录每个 goroutine 成功和失败的次数，最后累加生成结果，避免了锁争用造成的性能损失
		var successCountList = make([]int, routineNum)
		var failCountList = make([]int, routineNum)
		// 创建发送协程
		for i := 0; i < routineNum; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for {
					// 超时处理
					var email string
					var ok bool
					ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
					select {
					// 从 channel 里取出 email 发送
					case email, ok = <-emailsChan:
						if !ok {
							cancel()
							return // 如果 channel 关闭，说明所有 email 已经发送完毕，goroutine 退出
						}
						// 发送邮件
						err := sendEmail(email)
						if err == nil {
							successCountList[i]++
						} else {
							failCountList[i]++
						}
					case <-ctx.Done():
						logger.Log.Error(nil, "邮件发送超时", "email", email, "routine", i, "timeout", TimeOut)
						continue
					}
					cancel()
				}
			}(i)
		}
		wg.Wait() // 等待所有 goroutine 完成
		// 统计总的成功和失败次数
		for i := 0; i < routineNum; i++ {
			TotalSuccessCount += successCountList[i]
			TotalFailCount += failCountList[i]
		}
		// 记录发送历史
		history := model.WeeklyEmailSendingHistory{
			AdminID:      AdminID,
			Emails:       userEmails,
			TotalCount:   TotalSuccessCount + TotalFailCount,
			SuccessCount: TotalSuccessCount,
			FailCount:    TotalFailCount,
			TimeOut:      config.Configs.Smtp.WeeklyEmail.TimeOut,
			RoutineNum:   routineNum,
		}
		infrastructure.PostgresDb.Create(&history)
	}()
}

// GetWeeklyEmailSendingHistory 获取周报邮件发送历史
func GetWeeklyEmailSendingHistory(pageNum int, pageSize int) ([]model.WeeklyEmailSendingHistory, int64, error) {
	var histories []model.WeeklyEmailSendingHistory
	var total int64
	err := infrastructure.PostgresDb.
		Model(&model.WeeklyEmailSendingHistory{}).
		Count(&total).Order("created_at desc").
		Offset((pageNum - 1) * pageSize).
		Limit(pageSize).
		Find(&histories).Error
	return histories, total, err
}

// sendEmail 周报邮件发送函数
func sendEmail(email string) error {
	var subject = "Weekly Email"
	body, err := utils.GenerateEmailBody("weekly-email.html", nil)
	if err != nil {
		return err
	}
	msg := utils.GenerateHTMLMsg(email, config.Configs.Smtp.SmtpUser, config.Configs.Smtp.SmtpNickname, subject, body)
	return utils.SendEmail(email, subject, string(msg))
}
