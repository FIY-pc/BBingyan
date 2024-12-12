package service

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/infrastructure"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/es"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
	"gorm.io/gorm"
)

// CreatePost 创建文章
func CreatePost(dto dto.CreatePostDTO) error {
	post := model.Post{
		Title:  dto.Title,
		NodeID: dto.NodeID,
		UserID: dto.AuthorID,
	}

	// 创建文章
	result := infrastructure.PostgresDb.Model(&model.Post{}).Create(&post)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "error", result.Error)
		return result.Error
	}

	// 创建文章内容
	content := model.Content{
		PostID: post.ID,
		Text:   dto.Text,
	}
	if err := infrastructure.PostgresDb.Model(&model.Content{}).Create(&content).Error; err != nil {
		logger.Log.Error(nil, ModelError, "error", err)
		return err
	}
	// 创建 ES 索引
	err := es.IndexPost(post, content)
	if err != nil {
		logger.Log.Error(nil, "Failed to index post", "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

// UpdatePost 更新文章
func UpdatePost(dto dto.UpdatePostDTO) error {
	updatesPost := make(map[string]interface{})
	if dto.Title != "" {
		updatesPost["title"] = dto.Title
	}
	// 更新文章基本信息
	result := infrastructure.PostgresDb.Model(&model.Post{}).Where("id = ?", dto.ID).Updates(updatesPost)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "ID", dto.ID, "error", result.Error)
		return result.Error
	}

	// 更新内容
	updatesContent := make(map[string]interface{})
	if dto.Text != "" {
		updatesContent["Text"] = dto.Text
	}
	result = infrastructure.PostgresDb.Model(&model.Content{}).Where("post_id = ?", dto.ID).Updates(updatesContent)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "ID", dto.ID, "error", result.Error)
		return result.Error
	}
	logger.Log.Info(nil, Success, "ID", dto.ID)
	return nil
}

// DeletePost 删除文章
func DeletePost(postID uint) error {
	return infrastructure.PostgresDb.Transaction(func(tx *gorm.DB) error {
		// 删除内容
		if err := infrastructure.PostgresDb.Model(&model.Content{}).Where("post_id = ?", postID).Delete(&model.Content{}).Error; err != nil {
			logger.Log.Error(nil, ModelError, "ID", postID, "error", err)
			return err
		}
		// 删除文章
		result := infrastructure.PostgresDb.Model(&model.Post{}).Where("id = ?", postID).Delete(&model.Post{})
		if result.Error != nil {
			logger.Log.Error(nil, ModelError, "ID", postID, "error", result.Error)
			return result.Error
		}
		logger.Log.Info(nil, Success, "ID", postID)
		return nil
	})
}

// GetPostInfo 获取文章基本信息
func GetPostInfo(postID uint) (dto.PostDTO, error) {
	var targetPost model.Post
	result := infrastructure.PostgresDb.Model(&model.Post{}).Where("id = ?", postID).First(&targetPost)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "ID", postID, "error", result.Error)
		return dto.PostDTO{}, errors.New(ModelError)
	}

	// 将数据库模型转换为 DTO
	postDTO := dto.PostDTO{
		ID:        targetPost.ID,
		Title:     targetPost.Title,
		AuthorID:  targetPost.UserID,
		CreatedAt: targetPost.CreatedAt,
		UpdatedAt: targetPost.UpdatedAt,
	}
	logger.Log.Info(nil, Success, "ID", postID)
	return postDTO, nil
}

// GetPostWithContent 获取文章及其内容
func GetPostWithContent(postID uint) (dto.PostWithContentDTO, error) {
	// 获取文章基本信息
	var targetPost model.Post
	result := infrastructure.PostgresDb.Model(&model.Post{}).Where("id = ?", postID).First(&targetPost)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "ID", postID, "error", result.Error)
		return dto.PostWithContentDTO{}, errors.New(ModelError)
	}
	// 获取文章内容
	var content model.Content
	result = infrastructure.PostgresDb.Model(&model.Content{}).Where("post_id = ?", postID).First(&content)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "ID", postID, "error", result.Error)
		return dto.PostWithContentDTO{}, errors.New(ModelError)
	}

	// 将数据库模型转换为 DTO
	postWithContentDTO := dto.PostWithContentDTO{
		Post: dto.PostDTO{
			ID:        targetPost.ID,
			NodeID:    targetPost.NodeID,
			Title:     targetPost.Title,
			AuthorID:  targetPost.UserID,
			CreatedAt: targetPost.CreatedAt,
			UpdatedAt: targetPost.UpdatedAt,
		},
		Text: content.Text,
	}
	logger.Log.Info(nil, Success, "ID", postID)
	return postWithContentDTO, nil
}

// GetPostContent 获取文章内容
func GetPostContent(postID uint) (string, error) {
	var content model.Content
	result := infrastructure.PostgresDb.Model(&model.Content{}).Where("post_id = ?", postID).First(&content)
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "ID", postID, "error", result.Error)
		return "", errors.New(ModelError)
	}
	logger.Log.Info(nil, Success, "ID", postID)
	return content.Text, nil
}
