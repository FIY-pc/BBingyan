package utils

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/golang-jwt/jwt"
	"strings"
)

const (
	BearerSchema = "Bearer "
)

// JwtClaims JWT声明结构
type JwtClaims struct {
	UID     uint `json:"uid"`
	IsAdmin bool `json:"is_admin"`
	jwt.StandardClaims
}

func (c *JwtClaims) Valid() error {
	if c.UID == 0 {
		return errors.New("invalid user id")
	}
	return c.StandardClaims.Valid()
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*JwtClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, BearerSchema)

	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Configs.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
