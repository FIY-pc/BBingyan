package util

import (
	"errors"
	"fmt"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/util/Param"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"math/rand"
	"net/smtp"
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
	key := Param.CaptchaKey(email)
	result := rdb.Set(ctx, key, captcha, time.Second*time.Duration(config.Config.Captcha.Timeout))
	if result.Err() != nil {
		return result.Err()
	}
	if err := rdb.Close(); err != nil {
		return err
	}
	return nil
}

// SendCaptcha 发送验证码
func SendCaptcha(email string, captcha string) error {
	// 配置信息
	user := config.Config.Email.SmtpUser
	nickname := config.Config.Email.SmtpNickname
	password := config.Config.Email.SmtpPassword
	host := config.Config.Email.SmtpHost
	port := config.Config.Email.SmtpPort
	auth := smtp.PlainAuth("", user, password, host)
	contentType := "Content-Type: text/plain; charset=UTF-8"
	// 邮件标题
	subject := "BBingyan验证码"
	// 邮件主体
	body := "您的验证码为：" + captcha + "\n该验证码有效期为5分钟，请尽快验证。"
	// 邮件内容
	rawMsg := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s", email, nickname, user, subject, contentType, body)
	msg := []byte(rawMsg)
	// 发送邮件
	// 检查间隔
	if !CaptchaInterval(email) {
		return errors.New("send captcha too often")
	}

	from := user
	addr := fmt.Sprintf("%s:%s", host, port)
	_ = smtp.SendMail(addr, auth, from, []string{email}, msg)
	return nil
}

// GenerateCaptcha 产生随机验证码
func GenerateCaptcha() string {
	captcha := make([]byte, 6)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		captcha[i] = byte(rand.Intn(10) + '0')
	}
	return string(captcha)
}

// CaptchaInterval 检查验证码发送间隔是否冷却
func CaptchaInterval(email string) bool {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: config.Config.Redis.Password,
		DB:       config.Config.Redis.DB,
	})
	key := Param.CaptchaKey(email)
	// 如果redis中没有验证码条目，通过
	_, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return true
	}
	// 查询过期时间
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return false
	}
	if config.Config.Captcha.Timeout-ttl < config.Config.Captcha.Interval {
		return false
	}
	return true
}
