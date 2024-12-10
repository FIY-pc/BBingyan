package service

import (
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/utils/logger"
)

func CreateComment(commentDTO dto.CommentDTO) error {
	comment := model.Comment{
		UserID: commentDTO.UserID,
		PostID: commentDTO.PostID,
		Text:   commentDTO.Text,
	}
	if err := model.PostgresDb.Model(&model.Comment{}).Create(&comment).Error; err != nil {
		logger.Log.Error(nil, ModelError, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

func DeleteComment(commentID uint) error {
	if err := model.PostgresDb.Model(&model.Comment{}).Where("id = ?", commentID).Delete(&model.Comment{}).Error; err != nil {
		logger.Log.Error(nil, ModelError, "error", err)
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

func GetCommentByID(commentID uint) (dto.CommentDTO, error) {
	var commentDTO dto.CommentDTO
	var comment model.Comment
	if err := model.PostgresDb.Model(&model.Comment{}).Where("id = ?", commentID).First(&comment).Error; err != nil {
		logger.Log.Error(nil, ModelError, "error", err)
		return commentDTO, err
	}
	commentDTO = dto.CommentDTO{
		UserID: comment.UserID,
		PostID: comment.PostID,
		Text:   comment.Text,
	}
	logger.Log.Info(nil, Success)
	return commentDTO, nil
}

func GetCommentsByPostID(postID uint) ([]dto.CommentDTO, error) {
	var commentDTOs []dto.CommentDTO
	var comments []model.Comment
	if err := model.PostgresDb.Model(&model.Comment{}).Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		logger.Log.Error(nil, ModelError, "error", err)
		return commentDTOs, err
	}
	logger.Log.Info(nil, Success)
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, dto.CommentDTO{
			ID:     comment.ID,
			UserID: comment.UserID,
			PostID: comment.PostID,
			Text:   comment.Text,
		})
	}
	return commentDTOs, nil
}

func GetCommentsByUserID(userID uint) ([]dto.CommentDTO, error) {
	var commentDTOs []dto.CommentDTO
	var comments []model.Comment
	if err := model.PostgresDb.Model(&model.Comment{}).Where("user_id = ?", userID).Find(&comments).Error; err != nil {
		logger.Log.Error(nil, ModelError, "error", err)
		return commentDTOs, err
	}
	logger.Log.Info(nil, Success)
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, dto.CommentDTO{
			ID:     comment.ID,
			UserID: comment.UserID,
			PostID: comment.PostID,
			Text:   comment.Text,
		})
	}
	return commentDTOs, nil
}
