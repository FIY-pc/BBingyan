package model

import "errors"

func CreateComment(comment Comment) error {
	if postgresDb == nil {
		return errors.New("DB is nil")
	}
	postgresDb.Create(comment)
	return nil
}

func UpdateComment(comment Comment) error {
	if postgresDb == nil {
		return errors.New("DB is nil")
	}
	postgresDb.Where("id", comment.ID).Updates(comment)
	return nil
}

func GetCommentByID(id uint) (*Comment, error) {
	if postgresDb == nil {
		return nil, errors.New("DB is nil")
	}
	resultComment := Comment{}
	result := postgresDb.Where("id", id).First(&resultComment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &resultComment, nil
}

func DeleteCommentByID(id uint) error {
	if postgresDb == nil {
		return errors.New("DB is nil")
	}
	postgresDb.Where("id", id).Delete(&Comment{})
	return nil
}

func GetCommentByPage(articleId uint, page int, pageSize int) ([]Comment, error) {
	if postgresDb == nil {
		return nil, errors.New("DB is nil")
	}
	var comments []Comment
	postgresDb.Where("article_id", articleId).Offset((page - 1) * pageSize).Limit(pageSize).Find(&comments)
	return comments, nil
}

func GetArticleCommentCount(articleId uint) (int64, error) {
	if postgresDb == nil {
		return 0, errors.New("DB is nil")
	}
	var commentNum int64
	postgresDb.Model(Comment{}).Where("article_id", articleId).Count(&commentNum)
	return commentNum, nil
}

func GetUserCommentCount(userId uint) (int64, error) {
	if postgresDb == nil {
		return 0, errors.New("DB is nil")
	}
	var commentNum int64
	postgresDb.Model(Comment{}).Where("user_id", userId).Count(&commentNum)
	return commentNum, nil
}
