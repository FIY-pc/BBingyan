package main

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/router"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/FIY-pc/BBingyan/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	config.LoadConfig()
	logger.NewLogger()
	model.NewPostgres()
	model.NewRedisClient()
	service.InitAdmin()
	router.InitRouter(e)

	log.Fatal(e.Start(config.Configs.Server.Host + ":" + config.Configs.Server.Port))
}
