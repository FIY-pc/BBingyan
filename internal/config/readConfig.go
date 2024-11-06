package config

import (
	"encoding/json"
	"os"
	"strings"
)

var Config configStruct

type configStruct struct {
	Server   ServerConfig   `json:"server"`
	Postgres PostgresConfig `json:"postgres"`
	Redis    RedisConfig    `json:"redis"`
	Jwt      JwtConfig      `json:"jwt"`
	Email    EmailConfig    `json:"email"`
}

type ServerConfig struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type PostgresConfig struct {
	Dsn string `json:"dsn"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type JwtConfig struct {
	Secret string `json:"secret"`
	Expire int    `json:"expire"`
}

type EmailConfig struct {
	SmtpUser     string `json:"smtp_user"`
	SmtpPassword string `json:"smtp_password"`
	SmtpNickname string `json:"smtp_nickname"`
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     string `json:"smtp_port"`
}

type CaptchaConfig struct {
	Length int `json:"length"`
	Expire int `json:"expire"`
}

func InitConfig() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parts := strings.Split(dir, "BBingyan")
	path := parts[0] + "BBingyan/Config/default.json"

	ConfigFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(ConfigFile, &Config)
	if err != nil {
		panic(err)
	}
}
