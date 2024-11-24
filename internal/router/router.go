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
	// Public
	e.GET("/captcha", controller.GetCaptcha)
	e.POST("/tokens", controller.Login)
	e.POST("/users/register", controller.Register)

	// User
	e.GET("/users", controller.UserInfo)
	e.POST("/users/login", controller.Login)
	e.PUT("/users", controller.UserUpdate)
	e.DELETE("/users", controller.UserDelete)

	// follow
	e.POST("/follows", controller.Follow)
	e.DELETE("/follows", controller.Unfollow)
	e.GET("/follows/num", controller.GetFollowerNum)
	e.GET("/follows/isfollow", controller.IsFollowed)

	// Article
	e.GET("/articles", controller.ArticleInfo)
	e.POST("/articles", controller.ArticleCreate)
	e.PUT("/articles", controller.ArticleUpdate)
	e.DELETE("/articles", controller.ArticleDelete)

	// comment
	e.POST("/comments", controller.CommentCreate)
	e.DELETE("/comments", controller.CommentDelete)
}
