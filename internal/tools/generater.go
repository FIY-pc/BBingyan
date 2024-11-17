package tools

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"math/rand"
	"time"
)

func GenerateRandName() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.New(rand.NewSource(time.Now().UnixNano()))
	length := rand.Intn(config.Config.User.Nickname.RandMax-config.Config.User.Nickname.RandMin) + config.Config.User.Nickname.RandMin
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
