package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/labstack/echo/v4"
	"net/http"
)

func ArticleInfo(c echo.Context) error {

	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Info article success",
		Data: nil,
	})
}

func ArticleCreate(c echo.Context) error {
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Create article success",
		Data: nil,
	})
}

func ArticleUpdate(c echo.Context) error {
	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Update article success",
		Data: nil,
	})
}

func ArticleDelete(c echo.Context) error {

	return c.JSON(http.StatusOK, params.Common200Resp{
		Code: http.StatusOK,
		Msg:  "Delete article success",
		Data: nil,
	})
}
