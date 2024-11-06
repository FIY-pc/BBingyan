package util

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

// AuthCaptcha 验证验证码是否正确
func AuthCaptcha(email string, captcha string) bool {
	var result bool
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	// 检查验证码是否正确
	key := "captcha:" + email + ":" + captcha
	val, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		result = false
	} else if err != nil {
		result = false
	} else {
		result = val == captcha
	}
	err = rdb.Close()
	if err != nil {
		result = false
	}
	return result
}

// AddCaptcha 记录验证码
func AddCaptcha(email string, captcha string) error {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	// 记录验证码，设置验证码过期时间
	key := "captcha:" + email + ":" + captcha
	rdb.Set(ctx, key, captcha, time.Second*time.Duration(config.Config.Captcha.Timeout))
	if err := rdb.Close(); err != nil {
		return err
	}
	return nil
}
