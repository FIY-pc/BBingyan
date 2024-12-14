package config

import "time"

type Config struct {
	Server   *ServerConfig        `yaml:"server"`
	Postgres *PostgresConfig      `yaml:"postgres"`
	User     *UserConfig          `yaml:"user"`
	JWT      *JWTConfig           `yaml:"jwt"`
	Bcrypt   *BcryptConfig        `yaml:"bcrypt"`
	Log      *LogConfig           `yaml:"log"`
	Redis    *RedisConfig         `yaml:"redis"`
	ES       *ElasticsearchConfig `yaml:"elasticsearch"`
	Smtp     *SmtpConfig          `yaml:"smtp"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type PostgresConfig struct {
	Dsn string `yaml:"dsn"`
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

type ElasticsearchConfig struct {
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
}
