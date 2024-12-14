package model

import "time"

type WeeklyEmailSendingHistory struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AdminID      uint      `json:"admin_id"`                        // 进行发送周报操作的管理员ID
	Emails       []string  `json:"emails" gorm:"type:json"`         // 接收周报的用户的邮箱
	TotalCount   int       `json:"total_count"`                     // 总共发送的邮件数量
	SuccessCount int       `json:"success_count"`                   // 发送成功的邮件数量
	FailCount    int       `json:"fail_count"`                      // 发送失败的邮件数量
	TimeOut      string    `json:"time_out"`                        // 发送超时时间配置
	RoutineNum   int       `json:"routine_num"`                     // 发送邮件的协程数量
	CreateAt     time.Time `json:"create_at" gorm:"autoCreateTime"` // 记录创建时间
}
