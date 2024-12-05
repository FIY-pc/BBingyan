package config

import "time"

type configStruct struct {
	Server   ServerConfig        `json:"server" yaml:"server"`
	Postgres PostgresConfig      `json:"postgres" yaml:"postgres"`
	Redis    RedisConfig         `json:"redis" yaml:"redis"`
	ES       ElasticsearchConfig `json:"elasticsearch" yaml:"elasticsearch"`
	Log      LogConfig           `json:"log" yaml:"log"`
	Jwt      JwtConfig           `json:"jwt" yaml:"jwt"`
	Email    EmailConfig         `json:"email" yaml:"email"`
	Captcha  CaptchaConfig       `json:"captcha" yaml:"captcha"`
	User     UserConfig          `json:"user" yaml:"user"`
}

type ServerConfig struct {
	Port string `json:"port" yaml:"port"`
	Host string `json:"host" yaml:"host"`
	Env  string `json:"env" yaml:"env"`
}

type PostgresConfig struct {
	Dsn string `json:"dsn" yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

type ElasticsearchConfig struct {
	Addr string `json:"addr" yaml:"addr"`
}

type LogConfig struct {
	LogFile   string `json:"logfile" yaml:"logfile"`
	MaxSize   int    `json:"maxsize" yaml:"maxsize"`
	MaxAge    int    `json:"max_age" yaml:"max_age"`
	Compress  bool   `json:"compress" yaml:"compress"`
	LocalTime bool   `json:"localtime" yaml:"localtime"`
}

type JwtConfig struct {
	Secret string `json:"secret" yaml:"secret"`
	Expire int    `json:"expire" yaml:"expire"`
}

type EmailConfig struct {
	SmtpUser     string `json:"smtp_user" yaml:"smtp_user"`
	SmtpPassword string `json:"smtp_password" yaml:"smtp_password"`
	SmtpNickname string `json:"smtp_nickname" yaml:"smtp_nickname"`
	SmtpHost     string `json:"smtp_host" yaml:"smtp_host"`
	SmtpPort     string `json:"smtp_port" yaml:"smtp_port"`
}

type CaptchaConfig struct {
	Length   int           `json:"length" yaml:"length"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`
	Interval time.Duration `json:"interval" yaml:"interval"`
}

type UserConfig struct {
	Nickname NicknameConfig `json:"nickname" yaml:"nickname"`
	Admin    AdminConfig    `json:"admin" yaml:"admin"`
}

type NicknameConfig struct {
	RandMin   int `json:"rand_min" yaml:"rand_min"`
	RandMax   int `json:"rand_max" yaml:"rand_max"`
	Maxlength int `json:"max_length" yaml:"max_length"`
}

type AdminConfig struct {
	Email    string `json:"email" yaml:"email"`
	Password string `json:"password" yaml:"password"`
}
