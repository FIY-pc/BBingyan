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

func CreatePost(c echo.Context) error {
	// Bind and validate dto
	var createPostDTO dto.CreatePostDTO
	if err := context.BindAndValid(c, createPostDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Create
	err := service.CreatePost(createPostDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, params.Response{
		Success: true,
		Message: "Post created successfully",
	})
}

func UpdatePost(c echo.Context) error {
	// Bind and validate dto
	var updatePostDTO dto.UpdatePostDTO
	if err := context.BindAndValid(c, updatePostDTO); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Update
	err := service.UpdatePost(updatePostDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Post updated successfully",
	})
}

func DeletePost(c echo.Context) error {
	// Get ID
	sid := c.Param("id")
	uid, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Delete
	err = service.DeletePost(uint(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Post deleted successfully",
	})
}

func GetPostInfo(c echo.Context) error {
	// Get ID
	sid := c.Param("id")
	uid, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Get post
	post, err := service.GetPostInfo(uint(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Get Post info",
		Data:    post,
	})
}

func GetPostWithContent(c echo.Context) error {
	// Get ID
	sid := c.Param("id")
	uid, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Get post
	post, err := service.GetPostWithContent(uint(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Get Post with content",
		Data:    post,
	})
}

func GetPostContent(c echo.Context) error {
	// Get ID
	sid := c.Param("id")
	uid, err := strconv.Atoi(sid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Get content
	content, err := service.GetPostContent(uint(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "Get Post content",
		Data:    content,
	})
}
