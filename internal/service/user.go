package service

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/infrastructure"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/utils"
	"gorm.io/gorm"
)

const (
	Success           = "success"
	ModelError        = "model error"
	HashPasswordError = "hash password failed"
	NotFoundError     = "resource not found"
	ParamsMissing     = "params missing"
	ParamsInvalid     = "params invalid"
)

// InitAdmin 初始化管理员用户
func InitAdmin() {
	var existingAdmin model.User
	result := infrastructure.PostgresDb.Model(&model.User{}).Where("uid = ?", 0).First(&existingAdmin)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		err := createAdminUser()
		if err != nil {
			logger.Log.Error(nil, ModelError, "error", err)
		}
	}
	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "error", result.Error)
	}

	logger.Log.Info(nil, Success)
}

// createAdminUser 创建管理员用户
func createAdminUser() error {
	hashPassword, err := utils.HashPassword(config.Configs.User.InitAdmin.Password)
	if err != nil {
		logger.Log.Warn(nil, HashPasswordError, "err", err)
		return err
	}

	adminUser := model.User{
		UID:      0,
		Email:    config.Configs.User.InitAdmin.Email,
		Password: hashPassword,
		IsAdmin:  true,
	}

	if err := infrastructure.PostgresDb.Model(&model.User{}).Create(&adminUser).Error; err != nil {
		logger.Log.Warn(nil, ModelError, "error", err)
		return err
	}

	logger.Log.Info(nil, Success)
	return nil
}

// CreateUser 创建用户
func CreateUser(userDTO dto.UserCreateDTO) error {
	// 构造新用户模型
	hashPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		logger.Log.Warn(nil, "", "err", err)
		return err
	}
	user := model.User{
		Email:    userDTO.Email,
		Password: hashPassword,
		Nickname: userDTO.Nickname,
	}

	result := infrastructure.PostgresDb.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		logger.Log.Warn(nil, ModelError, "error", result.Error)
		return result.Error
	}

	logger.Log.Info(nil, Success, "UID", user.UID)
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(userDTO dto.UserUpdateDTO) error {
	// 构造更新字段
	updates := make(map[string]interface{})
	if userDTO.Nickname != "" {
		updates["nickname"] = userDTO.Nickname
	}
	if userDTO.Password != "" {
		hashPassword, err := utils.HashPassword(userDTO.Password)
		if err != nil {
			logger.Log.Warn(nil, HashPasswordError, "err", err)
			return err
		}
		updates["password"] = hashPassword
	}
	// 更新用户信息
	result := infrastructure.PostgresDb.Model(&model.User{}).Where("uid = ?", userDTO.UID).Updates(updates)
	if result.Error != nil {
		logger.Log.Warn(nil, ModelError, "UID", userDTO.UID, "error", result.Error)
		return result.Error
	}

	logger.Log.Info(nil, Success, "UID", userDTO.UID)
	return nil
}

// DeleteUser 删除用户
func DeleteUser(UID uint) error {
	result := infrastructure.PostgresDb.Model(&model.User{}).Where("uid = ?", UID).Delete(&model.User{})
	if result.Error != nil {
		logger.Log.Warn(nil, ModelError, "UID", UID)
		return result.Error
	}

	logger.Log.Info(nil, Success, "UID", UID)
	return nil
}

// GetUserByNickname 根据昵称获取用户
func GetUserByNickname(nickname string) (*model.User, error) {
	var user model.User
	result := infrastructure.PostgresDb.Model(&model.User{}).Where("nickname = ?", nickname).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.Warn(nil, NotFoundError, "nickname", nickname)
		return nil, errors.New(NotFoundError) // 用户未找到
	}

	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "nickname", nickname, "error", result.Error)
		return nil, result.Error
	}
	logger.Log.Info(nil, Success, "UID", user.UID)
	return &user, nil
}

// GetUserByEmail 根据电子邮件获取用户
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := infrastructure.PostgresDb.Model(&model.User{}).Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.Warn(nil, NotFoundError, "email", email)
		return nil, errors.New(NotFoundError) // 用户未找到
	}

	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "email", email, "error", result.Error)
		return nil, result.Error
	}
	logger.Log.Info(nil, Success, "UID", user.UID)
	return &user, nil
}

// GetUserByID 根据用户 ID 获取用户
func GetUserByID(uid uint) (*model.User, error) {
	var user model.User
	result := infrastructure.PostgresDb.Model(&model.User{}).Where("uid = ?", uid).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.Log.Warn(nil, NotFoundError, "UID", uid)
		return nil, errors.New(NotFoundError) // 用户未找到
	}

	if result.Error != nil {
		logger.Log.Error(nil, ModelError, "UID", uid, "error", result.Error)
		return nil, result.Error
	}
	logger.Log.Info(nil, Success, "UID", uid)
	return &user, nil
}
