package config

import (
	"os"
	"strings"
)

func DevGetConfigPath(configName string) string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	parts := strings.Split(dir, "BBingyan")
	path := parts[0] + "BBingyan/Config/" + configName
	return path
}
