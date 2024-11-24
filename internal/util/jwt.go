package util

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// JwtClaims 是一个结构体，用于存储JWT令牌的声明信息。
type JwtClaims struct {
	UserId     uint  `json:"userId"`
	Permission int   `json:"permission"`
	Exp        int64 `json:"exp"`
}

// Valid 方法用于验证 JwtClaims 结构体中的 token 是否有效。
func (c JwtClaims) Valid() error {
	if jwt.TimeFunc().Unix() > c.Exp {
		return jwt.NewValidationError("token is expired", jwt.ValidationErrorExpired)
	}
	return nil
}

// GenerateToken 用于根据传入的JwtClaim结构体生成token
func GenerateToken(claims JwtClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Config.Jwt.Secret))
}

// ParseToken 用于解析token
func ParseToken(tokenString string) (*JwtClaims, error) {
	// 去掉前缀Bearer，验证长度
	if len(tokenString) > 7 && tokenString[0:7] == "Bearer " {
		tokenString = tokenString[7:]
	} else {
		return &JwtClaims{}, jwt.NewValidationError("token is not a bearer token", jwt.ValidationErrorMalformed)
	}
	// 解析token成JwtClaim结构体
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Jwt.Secret), nil
	})
	if err != nil {
		return &JwtClaims{}, err
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return &JwtClaims{}, jwt.NewValidationError("invalid token", jwt.ValidationErrorMalformed)
}

// JWTAuthMiddleware 用于鉴权，包含token有效性验证和权限级别验证
func JWTAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 跳过不需要鉴权的路径
			if Skipper(c) {
				return next(c)
			}
			// 从请求中获取Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header missing")
			}
			// 解析并验证JWT令牌
			claims, err := ParseToken(authHeader)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			// 将解析后的claims存入上下文，供后续处理器使用
			c.Set("claims", claims)
			// 检查权限等级
			return PermissionMiddleware()(next)(c)
		}
	}
}

// PermissionMiddleware 权限级别验证
func PermissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := c.Get("claims").(*JwtClaims)
			permission := claims.Permission
			if _, exist := config.PathLevel[c.Path()][c.Request().Method]; !exist {
				return next(c)
			}
			if permission < config.PathLevel[c.Path()][c.Request().Method] {
				return echo.NewHTTPError(http.StatusUnauthorized, "permission denied")
			}
			return next(c)
		}
	}
}

func Skipper(c echo.Context) bool {
	if level, exist := config.PathLevel[c.Path()][c.Request().Method]; !exist || level != model.PermissionPublic {
		return false
	}
	return true
}
