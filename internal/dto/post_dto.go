package dto

import (
	"time"
)

// PostDTO 代表文章的基本信息
type PostDTO struct {
	ID        uint      `json:"id"`         // 文章ID
	NodeID    uint      `json:"node_id"`    // 文章节点ID
	Title     string    `json:"title"`      // 文章标题
	AuthorID  uint      `json:"author_id"`  // 作者ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// PostWithContentDTO 包含内容的文章
type PostWithContentDTO struct {
	Post PostDTO `json:"post"` // 文章基本信息
	Text string  `json:"text"` // 文章内容
}

// CreatePostDTO 用于创建文章的输入数据
type CreatePostDTO struct {
	Title    string `json:"title" validate:"required"`     // 文章标题
	AuthorID uint   `json:"author_id" validate:"required"` // 作者ID
	NodeID   uint   `json:"node_id" validate:"required"`   // 文章节点ID
	Text     string `json:"text"`                          // 文章内容
}

// UpdatePostDTO 用于更新文章的输入数据
type UpdatePostDTO struct {
	ID    uint   `json:"id" validate:"required"` // 文章ID
	Title string `json:"title"`                  // 文章标题
	Text  string `json:"Text"`                   // 文章内容
}
