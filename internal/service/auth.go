package service

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/infrastructure"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"math/rand"
	"strconv"
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
	result := infrastructure.PostgresDb.Model(&model.User{}).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	logger.Log.Info(c, Success)
	return nil
}

// SendCaptchaEmail 发送验证码邮件
func SendCaptchaEmail(email string) error {
	ctx := context.Background()
	// check interval
	err := checkInterval(ctx, email)
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return err
	}
	// send email
	captcha := generateCaptcha()
	subject := "BBingyan验证码"
	body, err := utils.GenerateEmailBody("captcha_email.html", captcha)
	if err != nil {
		logger.Log.Error(nil, err.Error())
		return err
	}
	_ = utils.SendEmail(email, subject, body)
	// add captcha to redis
	err = addCaptcha(ctx, email, captcha)
	if err != nil {
		return err
	}
	logger.Log.Info(nil, Success)
	return nil
}

func checkInterval(ctx context.Context, email string) error {
	rawTTL, err := time.ParseDuration(config.Configs.Smtp.Captcha.Expire)
	interval, err := time.ParseDuration(config.Configs.Smtp.Captcha.Interval)
	if err != nil {
		return err
	}
	TTL, err := infrastructure.Rdb.TTL(ctx, email).Result()
	if err != nil {
		return err
	}
	if TTL > 0 && rawTTL-TTL < interval {
		return errors.New("请勿频繁发送验证码")
	}
	return nil
}

func addCaptcha(ctx context.Context, email, captcha string) error {
	TTL, err := time.ParseDuration(config.Configs.Smtp.Captcha.Expire)
	if err != nil {
		return err
	}
	return infrastructure.Rdb.Set(ctx, email, captcha, TTL).Err()
}

func verifyCaptcha(email, captcha string) bool {
	key := email + ":" + captcha
	result, err := infrastructure.Rdb.Get(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if result != captcha {
		return false
	}
	_ = infrastructure.Rdb.Del(context.Background(), key)
	return true
}

func generateCaptcha() (captcha string) {
	return strconv.Itoa(100000 + rand.Intn(900000))
}
