package controller

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/FIY-pc/BBingyan/internal/tools"
	"github.com/FIY-pc/BBingyan/internal/util"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	captcha := c.FormValue("captcha")

	if email == "" || password == "" || captcha == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Login failed", errors.New("email,password and captcha are required"))
	}
	// 验证邮箱
	user, err := model.GetUserByEmail(email)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Login failed", err)
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Login failed", err)
	}
	// 验证验证码
	ok := util.AuthCaptcha(email, captcha)
	if !ok {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Login failed", err)
	}
	// 生成token
	claims := util.JwtClaims{
		UserId:     user.ID,
		Permission: user.Permission,
		Exp:        time.Now().Add(time.Minute * time.Duration(config.Config.Jwt.Expire)).Unix(),
	}
	token, err := util.GenerateToken(claims)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Login failed", err)
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Login success",
		Data: map[string]string{
			"token":    token,
			"nickname": user.Nickname,
		},
	})
}

func GetCaptcha(c echo.Context) error {
	var err error
	email := c.FormValue("email")
	// 发送验证码
	captcha := util.GenerateCaptcha()
	err = util.SendCaptcha(email, captcha)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "send captcha failed", err)
	}
	// 记录验证码
	err = util.AddCaptcha(email, captcha)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Add captcha failed", err)
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Send captcha success",
		Data: nil,
	})
}

func Register(c echo.Context) error {
	var err error
	email := c.FormValue("email")
	password := c.FormValue("password")
	captcha := c.FormValue("captcha")

	if email == "" || password == "" || captcha == "" {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Register failed", errors.New("email,password and captcha is required"))
	}
	// 检查该邮箱是否已注册
	_, err = model.GetUserByEmail(email)
	if err == nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Register failed", errors.New("email already exist"))
	}
	// 验证邮箱
	ok := util.AuthCaptcha(email, captcha)
	if !ok {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "Register failed", err)
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Register failed", err)
	}
	// 注册用户
	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Nickname: tools.GenerateRandName(),
	}
	err = model.CreateUser(user)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "Register failed", err)
	}
	// 返回成功响应
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Register success",
		Data: map[string]string{
			"email":    email,
			"nickname": user.Nickname,
		},
	})
}
