package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
)

func loadConfig() {
	env := os.Getenv("ENV")
	viper.AddConfigPath("./Config")
	if env != "" {
		viper.SetConfigName(env)
	} else {
		viper.SetConfigName("default")
	}
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		loadConfig()
	})
}

func InitConfigByViper() {
	// 开启从环境变量读取配置
	viper.AutomaticEnv()
	if err := viper.BindEnv("postgres.dsn", "POSTGRES_DSN"); err != nil {
		log.Fatalf("Unable to bind PostgreSQL DSN: %v", err)
	}
	if err := viper.BindEnv("redis.addr", "REDIS_ADDR"); err != nil {
		log.Fatalf("Unable to bind Redis ADDR: %v", err)
	}
	if err := viper.BindEnv("redis.password", "REDIS_PASSWORD"); err != nil {
		log.Fatalf("Unable to bind Redis PASSWORD: %v", err)
	}
	if err := viper.BindEnv("redis.db", "REDIS_DB"); err != nil {
		log.Fatalf("Unable to bind Redis DB: %v", err)
	}
	// 读取配置
	loadConfig()
	// 配置热更新
	watchConfig()
}
