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
		UserId:     user.ID,
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
	err = util.SendCaptcha(email, captcha)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
			Code:  http.StatusBadRequest,
			Msg:   "Sending captcha failed",
			Error: err.Error(),
		})
	}
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

func UserInfo(c echo.Context) error {
	email := c.FormValue("email")
	user := &model.User{}
	user, err := model.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get user failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User info success",
		Data: user,
	})
}

func UserUpdate(c echo.Context) error {
	email := c.FormValue("email")
	// 权限验证环节,仅自己和超级管理员可更改
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	permission := claims.Permission
	resultUser, err := model.GetUserByID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get user failed",
			Error: err.Error(),
		})
	}
	if permission < util.PermissionAdmin {
		if resultUser.Email != email {
			return c.JSON(http.StatusUnauthorized, params.CommonErrorResp{
				Code:  http.StatusUnauthorized,
				Msg:   "You are not allowed to delete this user",
				Error: "",
			})
		}
	}
	user, err := model.GetUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get user failed",
			Error: err.Error(),
		})
	}

	// 以下为比较适合在本路径更新的条目
	if Intro := c.FormValue("intro"); Intro != "" {
		user.Intro = Intro
	}
	if Password := c.FormValue("password"); Password != "" {
		user.Password = Password
	}
	if Nickname := c.FormValue("nickname"); Nickname != "" {
		if len(Nickname) > config.Config.User.Nickname.Maxlength {
			return c.JSON(http.StatusBadRequest, params.CommonErrorResp{
				Code:  http.StatusBadRequest,
				Msg:   "Nickname is too long",
				Error: "",
			})
		}
		user.Nickname = Nickname
	}

	err = model.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Update user failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User update success",
		Data: user,
	})
}

func UserDelete(c echo.Context) error {
	email := c.FormValue("email")
	claims := c.Get("claims").(util.JwtClaims)
	userId := claims.UserId
	permission := claims.Permission
	resultUser, err := model.GetUserByID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get user failed",
			Error: err.Error(),
		})
	}
	// 要么超级管理员删号，要么自己注销
	// TODO: 如果是用户自行注销，后续需要一个再次确认环节，需要正确验证码才能成功删号
	if permission < util.PermissionAdmin {
		if resultUser.Email != email {
			return c.JSON(http.StatusUnauthorized, params.CommonErrorResp{
				Code:  http.StatusUnauthorized,
				Msg:   "You are not allowed to delete this user",
				Error: "",
			})
		}
	}
	err = model.DeleteUserByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Delete user failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User delete success",
		Data: nil,
	})
}

// GetUserByNickName 根据昵称获取用户信息,主要为@功能提供支持
func GetUserByNickName(c echo.Context) error {
	nickname := c.QueryParam("nickname")
	user, err := model.GetUserByNickname(nickname)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.CommonErrorResp{
			Code:  http.StatusInternalServerError,
			Msg:   "Get user failed",
			Error: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User info success",
		Data: user,
	})
}
