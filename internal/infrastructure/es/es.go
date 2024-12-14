package es

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/elastic/go-elasticsearch/v7"
)

var ES *elasticsearch.Client

func NewElasticSearch() {
	var err error
	cfg := elasticsearch.Config{
		Addresses: config.Configs.ES.Addresses,
		Username:  config.Configs.ES.Username,
		Password:  config.Configs.ES.Password,
	}
	ES, err = elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Log.Fatal(nil, err.Error())
	}
	_, err = ES.Ping()
	if err != nil {
		ES = nil
		logger.Log.Fatal(nil, err.Error())
	}
}
