package model

import "time"

type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`             // 评论ID
	PostID    uint      `gorm:"not null" json:"post_id"`          // 关联的文章ID
	UserID    uint      `gorm:"not null" json:"author_id"`        // 评论者ID
	Text      string    `gorm:"type:text;not null" json:"text"`   // 评论内容
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // 创建时间
}
