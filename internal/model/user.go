package model

import (
	"time"
)

type User struct {
	UID                  uint      `json:"uid" gorm:"primaryKey;autoIncrement"`
	Email                string    `json:"email" gorm:"uniqueIndex;type:varchar(255);not null"`
	Password             string    `json:"-" gorm:"type:varchar(255);not null"` // 密码不输出到JSON
	IsAdmin              bool      `json:"is_admin" gorm:"default:false"`
	Nickname             string    `json:"nickname" gorm:"unique"`
	Intro                string    `json:"intro" gorm:"type:varchar(500);default:'无'"`
	Avatar               string    `json:"avatar" gorm:"type:varchar(255);default:'/static/avatar/default.png'"`
	SubscribeWeeklyEmail bool      `json:"subscribe_weekly_email" gorm:"default:false"` // 是否订阅周报
	CreateAt             time.Time `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt             time.Time `json:"update_at" gorm:"autoUpdateTime"`
}
