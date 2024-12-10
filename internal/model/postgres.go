package model

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var PostgresDb *gorm.DB

func NewPostgres() {
	var err error
	PostgresDb, err = gorm.Open(postgres.Open(config.Configs.Postgres.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// migrate user
	if err = PostgresDb.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}
	// migrate post related
	if err = PostgresDb.AutoMigrate(&Post{}); err != nil {
		log.Fatal(err)
	}
	if err = PostgresDb.AutoMigrate(&Content{}); err != nil {
		log.Fatal(err)
	}
	if err = PostgresDb.AutoMigrate(&Comment{}); err != nil {
		log.Fatal(err)
	}
	// migrate post and user related
	if err = PostgresDb.AutoMigrate(&Node{}); err != nil {
		log.Fatal(err)
	}
}
