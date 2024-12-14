package middleware

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/infrastructure"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// RateLimitMiddleware 时间间隔中间件，主要用于限制某些接口的访问频率，比如发送周报接口，防止误操作
func RateLimitMiddleware(rawTime string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			TTL, err := time.ParseDuration(rawTime)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, params.Response{
					Success: false,
					Message: "时间解析错误",
				})
			}
			path := c.Path()
			cmd := infrastructure.Rdb.Exists(c.Request().Context(), path)
			if cmd.Val() == 1 {
				return c.JSON(http.StatusTooManyRequests, params.Response{
					Success: false,
					Message: "操作过于频繁，请稍后再试",
				})
			}
			infrastructure.Rdb.Set(c.Request().Context(), path, 1, TTL)
			return next(c)
		}
	}
}
