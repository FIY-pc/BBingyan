package router

import (
	"github.com/FIY-pc/BBingyan/internal/controller"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
)

func InitRouter(e *echo.Echo) {
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
	e.GET("/users/commentCount", controller.GetUserCommentCount)

	// follow
	e.POST("/follows", controller.Follow)
	e.DELETE("/follows", controller.Unfollow)
	e.GET("/follows/isfollowed", controller.IsFollowed)
	e.GET("/followers/num", controller.GetFollowerNum)
	e.GET("/followers", controller.ListFollower)

	// Article
	e.GET("/articles", controller.ArticleInfo)
	e.POST("/articles", controller.ArticleCreate)
	e.PUT("/articles", controller.ArticleUpdate)
	e.DELETE("/articles", controller.ArticleDelete)
	e.GET("/article/commentsCount", controller.GetArticleCommentCount)
	// like
	e.POST("/articles/like", controller.Like)
	e.DELETE("/articles/unlike", controller.Unlike)
	e.GET("/articles/likeNum", controller.GetLikeNum)

	// comment
	e.POST("/comments", controller.CommentCreate)
	e.DELETE("/comments", controller.CommentDelete)
	e.GET("/comments", controller.CommentGetById)
	e.GET("/comments/pages", controller.CommentList)

	// node
	e.GET("/nodes", controller.NodeInfo)
	e.POST("/nodes", controller.CreateNode)
	e.DELETE("/nodes", controller.DeleteNode)
	e.PUT("/nodes", controller.UpdateNode)
	e.GET("/nodes/articles", controller.ListArticleFromNode)
	// node admin
	e.POST("/nodes/admins", controller.AddNodeAdmin)
	e.DELETE("/nodes/admins", controller.DeleteNodeAdmin)
	e.GET("/nodes/admins", controller.ListNodeAdmin)
}
