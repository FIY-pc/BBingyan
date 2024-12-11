package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	Configs *Config
)

func LoadConfig() {
	env := os.Getenv("ENV")
	viper.AddConfigPath(getConfigsPath())
	viper.AutomaticEnv()

	switch env {
	case "dev":
		viper.SetConfigName("config_dev")
	case "test":
		viper.SetConfigName("config_test")
	case "prod":
		viper.SetConfigName("config_prod")
	default:
		viper.SetConfigName("config_dev")
	}

	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config: " + err.Error())
	}

	if err := viper.Unmarshal(&Configs); err != nil {
		panic("Failed to unmarshal config: " + err.Error())
	}

	if Configs == nil {
		panic("Failed to load config")
	}
	startWatchingConfig()
}

// startWatchingConfig 监控配置文件变化并热重载
func startWatchingConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		var newConfig Config
		if err := viper.Unmarshal(&newConfig); err != nil {
			log.Fatal("Failed to unmarshal config on change: " + err.Error())
		} else {
			Configs = &newConfig
			log.Info("Configuration updated: " + e.Name)
		}
	})
}

func getConfigsPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	// 使用 filepath 包来处理路径
	path := strings.Split(dir, "BBingyan")[0]
	configPath := filepath.Join(path, "BBingyan/configs")

	return configPath
}
