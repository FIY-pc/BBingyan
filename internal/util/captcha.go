package util

import (
	"errors"
	"fmt"
	"github.com/FIY-pc/BBingyan/internal/config"
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
	key := "captcha:" + email + ":" + captcha
	rdb.Set(ctx, key, captcha, time.Second*time.Duration(config.Config.Captcha.Timeout))
	if err := rdb.Close(); err != nil {
		return err
	}
	return nil
}

func SendCaptcha(address string, captcha string) error {
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
	rawMsg := fmt.Sprintf("To:%s\r\nFrom:%s<%s>\r\nSubject:%s\r\n%s\r\n\r\n%s", address, nickname, user, subject, contentType, body)
	msg := []byte(rawMsg)
	// 发送邮件
	// TODO: 发送验证码要有间隔时间
	from := user
	addr := fmt.Sprintf("%s:%s", host, port)
	err := smtp.SendMail(addr, auth, from, []string{address}, msg)
	return err
}

func GenerateCaptcha() string {
	captcha := make([]byte, 6)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 6; i++ {
		captcha[i] = byte(rand.Intn(10) + '0')
	}
	return string(captcha)
}
