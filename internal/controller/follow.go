package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Follow(c echo.Context) error {
	followerUID, err := context.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	targetUID, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if followerUID == targetUID {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "You can't follow yourself",
		})
	}
	err = service.Follow(targetUID, followerUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Follow successfully",
	})
}

func UnFollow(c echo.Context) error {
	followerUID, err := context.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	targetUID, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if followerUID == targetUID {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: "You can't unfollow yourself",
		})
	}
	err = service.UnFollow(targetUID, followerUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Unfollow successfully",
	})
}
