package router

import (
	articlecontroller "github.com/FIY-pc/BBingyan/internal/controller/article"
	commentcontroller "github.com/FIY-pc/BBingyan/internal/controller/comment"
	followcontroller "github.com/FIY-pc/BBingyan/internal/controller/follow"
	likecontroller "github.com/FIY-pc/BBingyan/internal/controller/like"
	nodecontroller "github.com/FIY-pc/BBingyan/internal/controller/node"
	usercontroller "github.com/FIY-pc/BBingyan/internal/controller/user"
	"github.com/FIY-pc/BBingyan/internal/middleware"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
)

func InitRouter(e *echo.Echo) {
	// 中间件
	e.Use(util.JWTAuthMiddleware)
	e.Use(middleware.TraceMiddleware)
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, BBingyan!")
	})
	// Public
	e.GET("/captcha", usercontroller.GetCaptcha)
	e.POST("/tokens", usercontroller.Login)
	e.POST("/users/register", usercontroller.Register)

	// User
	e.GET("/users", usercontroller.UserInfo)
	e.PUT("/users", usercontroller.UserUpdate)
	e.DELETE("/users", usercontroller.UserDelete)
	e.GET("/users/commentCount", commentcontroller.GetUserCommentCount)

	// follow
	e.POST("/follows", followcontroller.Follow)
	e.DELETE("/follows", followcontroller.Unfollow)
	e.GET("/follows/isfollowed", followcontroller.IsFollowed)
	e.GET("/followers/num", followcontroller.GetFollowerNum)
	e.GET("/followers", followcontroller.ListFollower)

	// Article
	e.GET("/articles", articlecontroller.ArticleInfo)
	e.POST("/articles", articlecontroller.ArticleCreate)
	e.PUT("/articles", articlecontroller.ArticleUpdate)
	e.DELETE("/articles", articlecontroller.ArticleDelete)
	e.GET("/article/commentsCount", commentcontroller.GetArticleCommentCount)
	// like
	e.POST("/articles/like", likecontroller.Like)
	e.DELETE("/articles/unlike", likecontroller.Unlike)
	e.GET("/articles/likeNum", likecontroller.GetLikeNum)

	// comment
	e.POST("/comments", commentcontroller.CommentCreate)
	e.DELETE("/comments", commentcontroller.CommentDelete)
	e.GET("/comments", commentcontroller.CommentGetById)
	e.GET("/comments/pages", commentcontroller.CommentList)

	// node
	e.GET("/nodes", nodecontroller.NodeInfo)
	e.POST("/nodes", nodecontroller.CreateNode)
	e.DELETE("/nodes", nodecontroller.DeleteNode)
	e.PUT("/nodes", nodecontroller.UpdateNode)
	e.GET("/nodes/articles", nodecontroller.ListArticleFromNode)
	// node admin
	e.POST("/nodes/admins", nodecontroller.AddNodeAdmin)
	e.DELETE("/nodes/admins", nodecontroller.DeleteNodeAdmin)
	e.GET("/nodes/admins", nodecontroller.ListNodeAdmin)
}
