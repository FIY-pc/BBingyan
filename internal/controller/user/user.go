package controller

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/controller/permission"
	"github.com/FIY-pc/BBingyan/internal/model"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UserInfo(c echo.Context) error {
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
	user := &model.User{}
	user, err = model.GetUserByID(userId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "User info failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User info success",
		Data: user,
	})
}

func UserUpdate(c echo.Context) error {
	// 获取userId
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
	// 权限验证环节,仅自己和超级管理员可更改
	if !permission.UserPermissionCheck(c, userId) {
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "User info failed", errors.New("permission denied"))
	}
	// 获取user
	user, err := model.GetUserByID(userId)
	// 以下为比较适合在本路径更新的条目
	if Intro := c.FormValue("intro"); Intro != "" {
		user.Intro = Intro
	}
	if Password := c.FormValue("password"); Password != "" {
		user.Password = Password
	}
	if Nickname := c.FormValue("nickname"); Nickname != "" {
		if len(Nickname) > config.Config.User.Nickname.Maxlength {
			return params.CommonErrorGenerate(c, http.StatusBadRequest, "Nickname too long", errors.New("nickname too long"))
		}
		user.Nickname = Nickname
	}
	// 更新user
	err = model.UpdateUser(user)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "User info failed", err)
	}
	// 返回结果
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User update success",
		Data: user,
	})
}

func UserDelete(c echo.Context) error {
	// 获取userId
	userId, err := params.GetUserId(c)
	if err != nil {
		return err
	}
	// 权限验证环节,仅自己和超级管理员可更改
	if !permission.UserPermissionCheck(c, userId) {
		return params.CommonErrorGenerate(c, http.StatusUnauthorized, "User info failed", errors.New("permission denied"))
	}
	err = model.DeleteUserByID(userId)
	if err != nil {
		return params.CommonErrorGenerate(c, http.StatusInternalServerError, "User info failed", err)
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
		return params.CommonErrorGenerate(c, http.StatusBadRequest, "User info failed", err)
	}
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "User info success",
		Data: user,
	})
}
