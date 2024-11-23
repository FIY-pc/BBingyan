package main

import (
	"fmt"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/router"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	config.InitConfig()
	model.InitPostgres()
	util.InitRedis()
	model.InitSuperAdmin()
	router.InitRouter(e)

	startURL := fmt.Sprint(config.Config.Server.Host, ":", config.Config.Server.Port)
	e.Logger.Fatal(e.Start(startURL))
}
