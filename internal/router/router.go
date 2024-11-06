package router

import (
	"github.com/labstack/echo/v4"
)

func InitRouter(e *echo.Echo) {
	InitPublicRouter(e)
}

func InitPublicRouter(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})
}
