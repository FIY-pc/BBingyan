package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func CreateUser(c echo.Context) error {
	// Bind and validate dto
	var userDTO dto.UserCreateDTO
	if err := context.BindAndValid(c, &userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	if err := service.CreateUser(userDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, params.Response{
		Success: true,
		Message: "User created successfully",
	})
}

func UpdateUser(c echo.Context) error {
	// Bind and validate dto
	var userDTO dto.UserUpdateDTO
	if err := context.BindAndValid(c, &userDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Update user
	if err := service.UpdateUser(userDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "User updated successfully",
	})
}

func DeleteUser(c echo.Context) error {
	uid, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if err = service.DeleteUser(uid); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
	})
}

func GetUser(c echo.Context) error {
	id := c.QueryParam("id")
	nickname := c.QueryParam("nickname")
	email := c.QueryParam("email")

	// 根据提供的查询参数进行查询
	if id != "" {
		// 根据 ID 查询用户
		uid, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, params.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		user, err := service.GetUserByID(uint(uid))
		if err != nil {
			return c.JSON(http.StatusNotFound, params.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, params.Response{
			Success: true,
			Data:    user,
		})
	} else if nickname != "" {
		// 根据昵称查询用户
		user, err := service.GetUserByNickname(nickname)
		if err != nil {
			return c.JSON(http.StatusNotFound, params.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, params.Response{
			Success: true,
			Data:    user,
		})
	} else if email != "" {
		// 根据邮箱查询用户
		user, err := service.GetUserByEmail(email)
		if err != nil {
			return c.JSON(http.StatusNotFound, params.Response{
				Success: false,
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusOK, params.Response{
			Success: true,
			Data:    user,
		})
	}

	// 如果没有提供任何查询参数，返回错误
	return c.JSON(http.StatusBadRequest, params.Response{
		Success: false,
		Message: "No query parameters provided",
	})
}
