package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/dto"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SearchPost(c echo.Context) error {
	var searchPostDTO dto.SearchPostDTO
	err := context.BindAndValid(c, searchPostDTO)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	// Search posts
	posts, err := service.SearchPosts(searchPostDTO)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data: map[string]interface{}{
			"posts": posts,
		},
	})
}
