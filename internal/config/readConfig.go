package config

import (
	"encoding/json"
	"os"
	"strings"
)

var Config configStruct
var PathLevel map[string]map[string]string

// InitConfig 初始化配置结构体
func InitConfig() {
	InitConfigByViper()
	InitPathLevel()
}

func InitPathLevel() {
	path := devGetConfigPath("PathLevel.json")

	ConfigFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(ConfigFile, &PathLevel)
	if err != nil {
		panic(err)
	}
}

func devGetConfigPath(configName string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parts := strings.Split(dir, "BBingyan")
	path := parts[0] + "BBingyan/Config/" + configName
	return path
}
