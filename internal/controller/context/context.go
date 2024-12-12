package context

import (
	"errors"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/FIY-pc/BBingyan/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"strconv"
)

func BindAndValid(c echo.Context, req interface{}) (err error) {
	logger.Log.Debug(c, "BindAndValid", "req", req.(map[string]interface{}))
	err = c.Bind(req)
	if err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return err
	}
	return nil
}

func GetUIDFromParams(c echo.Context) (uint, error) {
	id := c.Param("id")
	if id == "" {
		return 0, errors.New("id is empty")
	}
	uid, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return uint(uid), nil
}

func GetUIDFromToken(c echo.Context) (uint, error) {
	claims := c.Get("claims").(utils.JwtClaims)
	uid := claims.UID
	if uid == 0 {
		return 0, errors.New("uid is empty")
	}
	return uid, nil
}
