package controller

import (
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
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Email, password and captcha are required",
			Error: "",
		})
	}
	// 验证邮箱
	user, err := model.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "User not found",
			Error: err.Error(),
		})
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Password is incorrect",
			Error: err.Error(),
		})
	}

	// 验证验证码
	ok := util.AuthCaptcha(email, captcha)
	if !ok {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Captcha is incorrect",
			Error: "",
		})
	}
	// 生成token
	claims := util.JwtClaims{
		ID:         user.ID,
		Email:      user.Email,
		Permission: user.Permission,
		Exp:        time.Now().Add(time.Minute * time.Duration(config.Config.Jwt.Expire)).Unix(),
	}
	token, err := util.GenerateToken(claims)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Generate token failed",
			Error: err.Error(),
		})
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
	_ = util.SendCaptcha(email, captcha)
	// 记录验证码
	err = util.AddCaptcha(email, captcha)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Add captcha failed",
			Error: err.Error(),
		})
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

	if email == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Email is null",
			Error: "",
		})
	}
	if password == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Password is null",
			Error: "",
		})
	}
	if captcha == "" {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Captcha is null",
			Error: "",
		})
	}
	// 检查该邮箱是否已注册
	_, err = model.GetUserByEmail(email)
	if err == nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Email already exists",
			Error: "",
		})
	}
	// 验证邮箱
	ok := util.AuthCaptcha(email, captcha)
	if !ok {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Captcha is incorrect",
			Error: "",
		})
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Hash password failed",
			Error: err.Error(),
		})
	}
	// 注册用户
	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Nickname: tools.GenerateRandName(),
	}
	err = model.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Create user failed",
			Error: err.Error(),
		})
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
