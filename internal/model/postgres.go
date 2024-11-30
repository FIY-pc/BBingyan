package model

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDb *gorm.DB

func InitPostgres() {
	var err error
	postgresDb, err = gorm.Open(postgres.Open(config.Config.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
	}
	InitUser(postgresDb)
	InitNode(postgresDb)
	InitArticle(postgresDb)
}
