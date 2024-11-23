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
