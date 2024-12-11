package service

import (
	"bytes"
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/utils"
	"github.com/FIY-pc/BBingyan/internal/utils/logger"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"html/template"
	"math/rand"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Login(dto dto.LoginDTO) (string, error) {
	var user *model.User
	var err error
	var token string
	// 查询账号
	if dto.Email != "" {
		user, err = GetUserByEmail(dto.Email)
		if err != nil {
			return "", err
		}
	} else {
		if dto.Nickname != "" {
			user, err = GetUserByNickname(dto.Nickname)
			if err != nil {
				return "", err
			}
		} else {
			logger.Log.Warn(nil, ParamsMissing)
			return "", errors.New(ParamsMissing)
		}
	}
	// 验证
	if !utils.ValidatePassword(user.Password, dto.Password) {
		logger.Log.Warn(nil, ParamsInvalid)
		return "", errors.New(ParamsInvalid)
	}
	// 生成 token
	token, err = generateToken(user.UID, user.IsAdmin)
	if err != nil {
		return "", err
	}
	logger.Log.Info(nil, Success)
	return token, nil
}

// generateToken 生成Token
func generateToken(uid uint, isAdmin bool) (string, error) {
	// 验证参数
	if uid == 0 {
		return "", errors.New("无效uid")
	}
	claims := &utils.JwtClaims{
		UID:     uid,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.Configs.JWT.Expiration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Configs.JWT.Secret))
}

func Register(c echo.Context, dto dto.RegisterDTO) error {
	var err error
	// 验证验证码
	if !verifyCaptcha(dto.Email, dto.Captcha) {
		return errors.New("验证码错误")
	}
	// 创建用户
	user := model.User{
		Email:    dto.Email,
		Nickname: dto.Nickname,
	}
	user.Password, err = utils.HashPassword(dto.Password)
	if err != nil {
		logger.Log.Warn(c, HashPasswordError)
		return err
	}
	result := model.PostgresDb.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	logger.Log.Info(c, Success)
	return nil
}

// SendCaptchaEmail 发送验证码邮件
func SendCaptchaEmail(email string) error {
	ctx := context.Background()
	// 检查发送间隔
	err := checkInterval(ctx, email)
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return err
	}
	// 生成验证码
	captcha := generateCaptcha()
	// 发送邮件
	subject := "BBingyan验证码"
	body, err := generateEmailBody(captcha)
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return err
	}
	msg := []byte("To: " + email + "\r\n" +
		"From: " + config.Configs.Captcha.SmtpUser + "\r\n" + "<" + config.Configs.Captcha.SmtpNickname + ">\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", config.Configs.Captcha.SmtpUser, config.Configs.Captcha.SmtpPassword, config.Configs.Captcha.SmtpHost)
	_ = smtp.SendMail(config.Configs.Captcha.SmtpHost+":"+config.Configs.Captcha.SmtpPort, auth, config.Configs.Captcha.SmtpUser, []string{email}, msg)
	// 将验证码存入redis
	err = addCaptcha(ctx, email, captcha)
	if err != nil {
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

func checkInterval(ctx context.Context, email string) error {
	rawTTL, err := time.ParseDuration(config.Configs.Captcha.Expire)
	interval, err := time.ParseDuration(config.Configs.Captcha.Interval)
	if err != nil {
		return err
	}
	TTL, err := model.Rdb.TTL(ctx, email).Result()
	if err != nil {
		return err
	}
	if TTL > 0 && rawTTL-TTL < interval {
		return errors.New("请勿频繁发送验证码")
	}
	return nil
}

func addCaptcha(ctx context.Context, email, captcha string) error {
	TTL, err := time.ParseDuration(config.Configs.Captcha.Expire)
	if err != nil {
		return err
	}
	return model.Rdb.Set(ctx, email, captcha, TTL).Err()
}

func verifyCaptcha(email, captcha string) bool {
	key := email + ":" + captcha
	result, err := model.Rdb.Get(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if result != captcha {
		return false
	}
	_ = model.Rdb.Del(context.Background(), key)
	return true
}

func generateCaptcha() (captcha string) {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

// generateEmailBody 从模板文件生成邮件内容
func generateEmailBody(captcha string) (string, error) {
	// 打开模板文件
	tmpl, err := template.ParseFiles(getTemplatePath())
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return "", err
	}

	// 创建一个字符串缓冲区来存储生成的内容
	var body bytes.Buffer
	err = tmpl.Execute(&body, captcha)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

func getTemplatePath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	// 使用 filepath 包来处理路径
	path := strings.Split(dir, "BBingyan")[0]
	templatePath := filepath.Join(path, "BBingyan"+"/web/templates/captcha_email.html")
	return templatePath
}
