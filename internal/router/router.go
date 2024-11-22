package router

import (
	"github.com/FIY-pc/BBingyan/internal/controller"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
)

func InitRouter(e *echo.Echo) {
	InitPublicRouter(e)
}

func InitPublicRouter(e *echo.Echo) {
	e.Use(util.JWTAuthMiddleware())
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, BBingyan!")
	})
	e.GET("/captcha", controller.GetCaptcha)
	e.POST("/tokens", controller.Login)
	e.POST("/users/register", controller.Register)
}
