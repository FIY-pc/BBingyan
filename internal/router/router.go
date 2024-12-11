package router

import (
	"github.com/FIY-pc/BBingyan/internal/controller"
	"github.com/FIY-pc/BBingyan/internal/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// InitRouter 初始化路由
func InitRouter(e *echo.Echo) {
	// common middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(middleware.TraceMiddleware)
	// route groups
	AuthRouter(e)
	UserRouter(e)
	PostRouter(e)
	NodeRouter(e)
	FollowRouter(e)
	CommentRouter(e)
}

// AuthRouter 认证路由
func AuthRouter(e *echo.Echo) {
	public := e.Group("/api/auth")
	public.POST("/login", controller.Login)
	public.POST("/captcha", controller.SendCaptcha)
	public.POST("/register", controller.Register)
}

// UserRouter 用户路由
func UserRouter(e *echo.Echo) {
	userGroup := e.Group("/api/users")

	// basic auth routes
	basic := userGroup.Group("")
	basic.Use(middleware.BasicAuth)
	basic.GET("/", controller.GetUser)
	basic.GET("/:id/comments", controller.GetCommentsByUserID)
	// owner auth routes
	owner := userGroup.Group("")
	owner.Use(middleware.OwnerAuth("user"))
	owner.PUT("/:id", controller.UpdateUser)
	owner.DELETE("/:id", controller.DeleteUser)

	// admin auth routes
	admin := userGroup.Group("")
	admin.Use(middleware.AdminAuth)
	admin.POST("/", controller.CreateUser)
}

// PostRouter 文章路由
func PostRouter(e *echo.Echo) {
	postGroup := e.Group("/api/posts")

	// basic auth routes
	basic := postGroup.Group("")
	basic.Use(middleware.BasicAuth)
	basic.GET("/", controller.GetPostWithContent)
	basic.GET("/info/:id", controller.GetPostInfo)
	basic.GET("/content/:id", controller.GetPostContent)
	basic.GET("/:id/comments", controller.GetCommentsByPostID)
	basic.POST("/", controller.CreatePost)
	// owner auth routes
	owner := postGroup.Group("")
	owner.Use(middleware.OwnerAuth("post"))
	owner.PUT("/:id", controller.UpdatePost)
	owner.DELETE("/:id", controller.DeletePost)
}

func NodeRouter(e *echo.Echo) {
	nodeGroup := e.Group("/api/nodes")

	basic := nodeGroup.Group("")
	basic.Use(middleware.BasicAuth)
	basic.GET("/:id", controller.GetNodeByID)
	basic.POST("/", controller.CreateNode)

	admin := nodeGroup.Group("")
	admin.Use(middleware.AdminAuth)
	admin.PUT("/:id", controller.UpdateNode)
	admin.DELETE("/:id/sort", controller.SortDeleteNode)
	admin.DELETE("/:id/hard", controller.HardDeleteNode)
	admin.DELETE("/posts", controller.DeletePostsUnderNode)
}

func FollowRouter(e *echo.Echo) {
	followGroup := e.Group("/api/follows")

	basic := followGroup.Group("")
	basic.Use(middleware.BasicAuth)
	basic.POST("/", controller.Follow)
	basic.DELETE("/", controller.UnFollow)
}

func CommentRouter(e *echo.Echo) {
	commentGroup := e.Group("/api/comments")

	basic := commentGroup.Group("")
	basic.Use(middleware.BasicAuth)
	basic.POST("/", controller.CreateComment)
	basic.GET("/:id", controller.GetCommentByID)

	owner := commentGroup.Group("")
	owner.Use(middleware.OwnerAuth("comment"))
	owner.DELETE("/:id", controller.DeleteComment)
}
