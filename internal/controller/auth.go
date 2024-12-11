package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/FIY-pc/BBingyan/internal/utils/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Register(c echo.Context) error {
	// Bind and validate dto
	var registerDTO dto.RegisterDTO
	if err := context.BindAndValid(c, registerDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Register
	err := service.Register(c, registerDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, params.Response{
		Success: true,
		Message: "Register successfully",
	})
}

func Login(c echo.Context) error {
	// Bind and validate dto
	var loginDTO dto.LoginDTO
	if err := context.BindAndValid(c, loginDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Login
	token, err := service.Login(loginDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Login successfully",
		Data: map[string]string{
			"token": token,
		},
	})
}

func SendCaptcha(c echo.Context) error {
	// Get email
	email := c.FormValue("email")
	if email == "" {
		logger.Log.Warn(c, "Email is required")
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "Email is required",
		})
	}
	// Send captcha email
	err := service.SendCaptchaEmail(email)
	if err != nil {
		logger.Log.Error(c, err.Error())
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	logger.Log.Info(c, "Send captcha successfully")
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Send captcha successfully",
	})
}
