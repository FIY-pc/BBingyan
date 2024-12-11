package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetCommentByID(c echo.Context) error {
	commentID, err := context.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	comment, err := service.GetCommentByID(commentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data:    comment,
	})
}

func GetCommentsByPostID(c echo.Context) error {
	postID, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	comments, err := service.GetCommentsByPostID(postID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data:    comments,
	})
}

func GetCommentsByUserID(c echo.Context) error {
	userID, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	comments, err := service.GetCommentsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data:    comments,
	})
}

func CreateComment(c echo.Context) error {
	var commentDTO dto.CommentDTO
	if err := context.BindAndValid(c, commentDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	err := service.CreateComment(commentDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "评论成功",
	})
}

func DeleteComment(c echo.Context) error {
	commentID, err := context.GetUIDFromParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	err = service.DeleteComment(commentID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "删除成功",
	})
}
