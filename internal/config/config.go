package config

import "time"

type Config struct {
	Server   *ServerConfig   `yaml:"server"`
	Postgres *PostgresConfig `yaml:"postgres"`
	User     *UserConfig     `yaml:"user"`
	Captcha  *CaptchaConfig  `yaml:"captcha"`
	JWT      *JWTConfig      `yaml:"jwt"`
	Bcrypt   *BcryptConfig   `yaml:"bcrypt"`
	Log      *LogConfig      `yaml:"log"`
	Redis    *RedisConfig    `yaml:"redis"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type UserConfig struct {
	InitAdmin *InitAdminConfig `yaml:"initAdmin"`
}

type PostgresConfig struct {
	Dsn string `yaml:"dsn"`
}

type InitAdminConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type CaptchaConfig struct {
	SmtpUser     string `yaml:"smtpUser"`
	SmtpNickname string `yaml:"smtpNickname"`
	SmtpPassword string `yaml:"smtpPassword"`
	SmtpHost     string `yaml:"smtpHost"`
	SmtpPort     string `yaml:"smtpPort"`
	Expire       string `yaml:"expire"`
	Interval     string `yaml:"interval"`
}

type JWTConfig struct {
	Secret     string        `yaml:"secret"`
	Expiration time.Duration `yaml:"expiration"`
}

type LogConfig struct {
	LogFile   string `yaml:"logfile"`
	MaxSize   int    `yaml:"maxsize"`
	MaxAge    int    `yaml:"maxAge"`
	Compress  bool   `yaml:"compress"`
	LocalTime bool   `yaml:"localTime"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type BcryptConfig struct {
	Cost int `yaml:"cost"`
}
