package model

import (
	"time"
)

// Post 代表文章模型
type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`             // 文章ID
	NodeID    uint      `gorm:"not null" json:"nid"`              // 文章节点ID
	Title     string    `gorm:"size:255;not null" json:"title"`   // 文章标题
	UserID    uint      `gorm:"not null" json:"user_id"`          // 作者ID
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // 更新时间
}

// Content 代表文章内容模型
type Content struct {
	ID     uint   `gorm:"primaryKey" json:"id"`           // 内容ID
	PostID uint   `gorm:"not null" json:"post_id"`        // 关联的文章ID
	Text   string `gorm:"type:text;not null" json:"text"` // 文章内容
}
