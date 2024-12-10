package utils

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), config.Configs.Bcrypt.Cost)
	return string(hashPassword), err
}

// ValidatePassword 验证密码
func ValidatePassword(hashPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}
