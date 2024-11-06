package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	PermissionPublic = 0
	PermissionUser   = 1
	PermissionAdmin  = 2
)

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Email      string    `json:"email" gorm:"unique"`
	Nickname   string    `json:"nickname"`
	Password   string    `json:"password"`
	Permission uint      `json:"permission"` // 权限级别，普通用户为0，节点管理员为1，最高管理员为2
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt"`

	UserActivity []UserActivity `json:"user_activity"`                      //has many
	UserInfo     UserInfo       `json:"user_info"`                          //has one
	UserHistory  []UserHistory  `json:"user_history"`                       //has many
	Article      []Article      `json:"article"`                            //has many
	Follower     []User         `json:"follower" gorm:"many2many:follows;"` //many to many
}

type UserActivity struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `json:"user_id"`
	ArticleNum  uint      `json:"article_num"`
	CommentNum  uint      `json:"comment_num"`
	LikeNum     uint      `json:"like_num"`
	FollowerNum uint      `json:"follower_num"`
	FocusNum    uint      `json:"focus_num"`
}

type UserInfo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	Intro     string    `json:"intro"`
	Avatar    string    `json:"avatar"`
}

type UserHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	ArticleID uint      `json:"article_id"`
	IsLiked   bool      `json:"is_liked"`
}

type UserFollower struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     uint      `json:"user_id"`
	FollowerID uint      `json:"follower_id"`
}

func InitUser(db *gorm.DB) {
	if err := db.AutoMigrate(&User{}, &UserFollower{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&UserActivity{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&UserInfo{}); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&UserHistory{}); err != nil {
		panic(err)
	}
}
