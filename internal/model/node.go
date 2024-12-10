package model

import (
	"gorm.io/gorm"
	"time"
)

type Node struct {
	ID        uint           `gorm:"primaryKey" json:"id"`                                                    // 节点ID
	Name      string         `gorm:"size:255;not null" json:"name"`                                           // 节点名称
	Intro     string         `gorm:"size:255;not null" json:"intro"`                                          // 节点简介
	Avatar    string         `gorm:"size:255;not null;default:/static/avatar/node/default.png" json:"avatar"` // 节点头像
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`                                        // 创建时间
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                                        // 更新时间
	DeleteAt  gorm.DeletedAt `gorm:"index" json:"delete_at"`                                                  // 删除时间, 软删除处理
}
