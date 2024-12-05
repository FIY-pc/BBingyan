package util

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// pathLevelJsonParser 用于建立PathLevel.json文件中字符串到权限数值的映射
var pathLevelJsonParser map[string]int

// InitPathLevelJsonParser  初始化用于解析PathLevel的映射
func InitPathLevelJsonParser() {
	pathLevelJsonParser = map[string]int{
		"public": PermissionPublic,
		"user":   PermissionUser,
		"admin":  PermissionAdmin,
	}
}

// 权限等级数字常量
const (
	PermissionPublic = 0
	PermissionUser   = 1
	PermissionAdmin  = 5
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
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

// PermissionMiddleware 权限级别验证,默认从PathLevel读取特殊路由权限配置
func PermissionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var level int
			var exist bool
			claims := c.Get("claims").(*JwtClaims)
			permission := claims.Permission
			if level, exist = getPermission(c); !exist {
				return next(c)
			}
			if permission < level {
				return echo.NewHTTPError(http.StatusUnauthorized, "permission denied")
			}
			return next(c)
		}
	}
}

// Skipper 用于跳过不需鉴权的路径
func Skipper(c echo.Context) bool {
	if level, exist := getPermission(c); !exist || level != PermissionPublic {
		return false
	}
	return true
}

// getPermission 获取请求路径的权限等级
func getPermission(c echo.Context) (int, bool) {
	result, exist := pathLevelJsonParser[config.PathLevel[c.Path()][c.Request().Method]]
	return result, exist
}
