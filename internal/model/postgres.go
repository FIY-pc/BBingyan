package model

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresDb *gorm.DB

func InitPostgres() {
	var err error
	postgresDb, err = gorm.Open(postgres.Open(config.Config.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	InitUser(postgresDb)
	InitArticle(postgresDb)
	InitNode(postgresDb)
}
