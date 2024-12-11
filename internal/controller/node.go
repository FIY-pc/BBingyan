package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetNodeByID(c echo.Context) error {
	nodeId, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	node, err := service.GetNodeByID(nodeId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data:    node,
	})
}

func CreateNode(c echo.Context) error {
	var nodeDTO dto.NodeDTO
	if err := context.BindAndValid(c, &nodeDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if err := service.CreateNode(nodeDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data:    nodeDTO,
	})
}

func UpdateNode(c echo.Context) error {
	var nodeDTO dto.NodeDTO
	if err := context.BindAndValid(c, &nodeDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if err := service.UpdateNode(nodeDTO); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data:    nodeDTO,
	})
}

func SortDeleteNode(c echo.Context) error {
	nodeId, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if err = service.SortDeleteNode(nodeId); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
	})
}

func HardDeleteNode(c echo.Context) error {
	nodeId, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	if err = service.HardDeleteNode(nodeId); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
	})
}

func DeletePostsUnderNode(c echo.Context) error {
	nodeId, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// delete all posts under this node
	if err = service.DeletePostUnderNode(nodeId); err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
	})
}
