package middleware

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/FIY-pc/BBingyan/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

const (
	ContextClaimKey = "claims"
)

// BasicAuth 基础认证中间件
func BasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		auth := c.Request().Header.Get("Authorization")
		if auth == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "无认证信息",
			})
		}

		claims, err := utils.ParseToken(auth)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "无效的token",
			})
		}

		c.Set(ContextClaimKey, claims)
		return next(c)
	}
}

// AdminAuth 管理员权限中间件
func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, ok := c.Get(ContextClaimKey).(*utils.JwtClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "未登录",
			})
		}

		if !claims.IsAdmin {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "需要管理员权限",
			})
		}

		return next(c)
	}
}

// OwnerAuth 资源所有者验证中间件，包含了管理员权限
func OwnerAuth(resourceType string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(ContextClaimKey).(*utils.JwtClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, params.Response{
					Success: false,
					Message: "unauthorized",
				})
			}

			if claims.IsAdmin {
				return next(c)
			}

			resourceID := c.Param("id")
			if resourceID == "" {
				return c.JSON(http.StatusBadRequest, params.Response{
					Success: false,
					Message: "invalid resource id",
				})
			}
			// 资源所有权验证逻辑
			switch resourceType {
			case "user":
				if err := userOwnerAuth(c, claims); err != nil {
					return err
				}
			case "post":
				if err := postOwnerAuth(c, claims); err != nil {
					return err
				}
			case "comment":
				if err := commentOwnerAuth(c, claims); err != nil {
					return err
				}
			}
			return next(c)
		}
	}
}

func postOwnerAuth(c echo.Context, claims *utils.JwtClaims) error {
	// 获取用户ID
	uid := claims.UID
	postIDStr := c.Param("id")
	if postIDStr == "" {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid user id",
		})
	}
	// 获取帖子ID
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid user id",
		})
	}
	// 获取帖子信息
	info, err := service.GetPostInfo(uint(postID))
	if err != nil {
		return c.JSON(http.StatusForbidden, params.Response{
			Success: false,
			Message: "forbidden",
		})
	}
	// 判断是否是帖子作者
	if info.AuthorID != uid {
		return c.JSON(http.StatusForbidden, params.Response{
			Success: false,
			Message: "forbidden",
		})
	}
	return nil
}

func userOwnerAuth(c echo.Context, claims *utils.JwtClaims) error {
	// 获取用户ID
	uid := claims.UID
	// 获取拥有者ID
	OwnerUIDStr := c.Param("id")
	if OwnerUIDStr == "" {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid user id",
		})
	}
	ownerUID, err := strconv.Atoi(OwnerUIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid user id",
		})
	}
	// 判断是否是用户本人
	if uid != uint(ownerUID) {
		return c.JSON(http.StatusForbidden, params.Response{
			Success: false,
			Message: "forbidden",
		})
	}
	return nil
}

func commentOwnerAuth(c echo.Context, claims *utils.JwtClaims) error {
	uid := claims.UID
	commentIDStr := c.Param("id")
	if commentIDStr == "" {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid comment id",
		})
	}
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid comment id",
		})
	}
	comment, err := service.GetCommentByID(uint(commentID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "invalid comment id",
		})
	}
	if comment.UserID != uid {
		return c.JSON(http.StatusForbidden, params.Response{
			Success: false,
			Message: "forbidden",
		})
	}
	return nil
}
