package model

import (
	"errors"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"time"
)

// Article 文章模型
type Article struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"` // 文章创建时间
	UpdatedAt time.Time `json:"updated_at"` // 文章更新时间
	UserID    uint      `json:"user_id"`    // 文章作者的ID
	Title     string    `json:"title"`      // 文章标题
	NodeID    uint      `json:"node_id"`
	Comment   []Comment `json:"comment" gorm:"foreignKey:ArticleID;-"` // forbid preload
	Content   Content   `json:"content" gorm:"foreignKey:ArticleID"`   // has one
}

// Content 文章内容表
type Content struct {
	ID        uint   `json:"id" gorm:"primarykey" `
	ArticleID uint   `json:"article_id"`
	Text      string `json:"text" gorm:"type:text"`
	UpdatedAt int64  `json:"updated_at"`
}

// Comment 文章评论表
type Comment struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	ArticleID uint      `json:"article_id"` // 文章ID
	UserID    uint      `json:"user_id"`    // 评论用户的ID
	Content   string    `json:"content"`    // 评论内容
	CreatedAt time.Time `json:"created_at"` // 评论时间
}

func InitArticle(db *gorm.DB) {
	if err := db.AutoMigrate(&Article{}); err != nil {
		log.Error(err)
	}
	if err := db.AutoMigrate(&Content{}); err != nil {
		log.Error(err)
	}
	if err := db.AutoMigrate(&Comment{}); err != nil {
		log.Error(err)
	}
}

func CreateArticle(article Article) error {
	if postgresDb == nil {
		return errors.New("DB is nil")
	}
	postgresDb.Create(article)
	return nil
}

func UpdateArticle(article Article) error {
	if postgresDb == nil {
		return errors.New("DB is nil")
	}
	postgresDb.Where("id", article.ID).Updates(article)
	return nil
}

func GetArticleByID(id uint) (*Article, error) {
	if postgresDb == nil {
		return nil, errors.New("DB is nil")
	}
	resultArticle := Article{}
	result := postgresDb.Where("id", id).First(&resultArticle)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resultArticle, nil
}

func DeleteArticleByID(ID uint) error {
	if postgresDb == nil {
		return errors.New("DB is nil")
	}
	postgresDb.Where("id", ID).Delete(&Article{})
	return nil
}
