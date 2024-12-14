package router

import (
	"github.com/FIY-pc/BBingyan/internal/config"
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
	authRouter(e)
	userRouter(e)
	postRouter(e)
	nodeRouter(e)
	followRouter(e)
	commentRouter(e)
	searchRouter(e)
	weeklyEmailRouter(e)
}

// authRouter 认证路由
func authRouter(e *echo.Echo) {
	public := e.Group("/api/auth")
	public.POST("/login", controller.Login)
	public.POST("/captcha", controller.SendCaptcha)
	public.POST("/register", controller.Register)
}

// userRouter 用户路由
func userRouter(e *echo.Echo) {
	// basic auth routes
	basic := e.Group("/api/users")
	basic.Use(middleware.BasicAuth)
	basic.GET("/", controller.GetUser)
	basic.GET("/:id/comments", controller.GetCommentsByUserID)
	// owner auth routes
	owner := basic.Group("")
	owner.Use(middleware.OwnerAuth("user"))
	owner.PUT("/:id", controller.UpdateUser)
	owner.DELETE("/:id", controller.DeleteUser)

	// admin auth routes
	admin := basic.Group("")
	admin.Use(middleware.AdminAuth)
	admin.POST("/", controller.CreateUser)
}

// postRouter 文章路由
func postRouter(e *echo.Echo) {
	// basic auth routes
	basic := e.Group("/api/posts")
	basic.Use(middleware.BasicAuth)
	basic.GET("/", controller.GetPostWithContent)
	basic.GET("/info/:id", controller.GetPostInfo)
	basic.GET("/content/:id", controller.GetPostContent)
	basic.GET("/:id/comments", controller.GetCommentsByPostID)
	basic.POST("/", controller.CreatePost)
	// owner auth routes
	owner := basic.Group("")
	owner.Use(middleware.OwnerAuth("post"))
	owner.PUT("/:id", controller.UpdatePost)
	owner.DELETE("/:id", controller.DeletePost)
}

func nodeRouter(e *echo.Echo) {
	basic := e.Group("/api/nodes")
	basic.Use(middleware.BasicAuth)
	basic.GET("/:id", controller.GetNodeByID)
	basic.POST("/", controller.CreateNode)

	admin := basic.Group("")
	admin.Use(middleware.AdminAuth)
	admin.PUT("/:id", controller.UpdateNode)
	admin.DELETE("/:id/sort", controller.SortDeleteNode)
	admin.DELETE("/:id/hard", controller.HardDeleteNode)
	admin.DELETE("/posts", controller.DeletePostsUnderNode)
}

func followRouter(e *echo.Echo) {
	basic := e.Group("/api/follows")
	basic.Use(middleware.BasicAuth)
	basic.POST("/", controller.Follow)
	basic.DELETE("/", controller.UnFollow)
}

func commentRouter(e *echo.Echo) {
	basic := e.Group("/api/comments")
	basic.Use(middleware.BasicAuth)
	basic.POST("/", controller.CreateComment)
	basic.GET("/:id", controller.GetCommentByID)

	owner := basic.Group("")
	owner.Use(middleware.OwnerAuth("comment"))
	owner.DELETE("/:id", controller.DeleteComment)
}

func searchRouter(e *echo.Echo) {
	basic := e.Group("/api/search")
	basic.Use(middleware.BasicAuth)
	basic.GET("/post", controller.SearchPost)
}

func weeklyEmailRouter(e *echo.Echo) {
	basic := e.Group("/api/weekly-email")
	basic.Use(middleware.BasicAuth)
	basic.POST("/subscribe", controller.SubscribeWeeklyEmail)

	admin := basic.Group("")
	admin.Use(middleware.AdminAuth)
	admin.GET("/history", controller.GetWeeklyEmailSendingHistory)
	admin.POST("/send", controller.SendWeeklyEmail, middleware.RateLimitMiddleware(config.Configs.Smtp.WeeklyEmail.RateLimit))
}
