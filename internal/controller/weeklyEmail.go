package controller

import (
	"github.com/FIY-pc/BBingyan/internal/controller/context"
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SendWeeklyEmail(c echo.Context) error {
	AdminID, err := context.GetUIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	service.SendWeeklyEmail(AdminID)
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Message: "后台发送进程已启动",
	})
}

// GetWeeklyEmailSendingHistory 获取周报邮件发送历史记录
func GetWeeklyEmailSendingHistory(c echo.Context) error {
	// 因为过于简单，所以不单独把这个req放到集中的地方了
	req := struct {
		Page     int `json:"page" validate:"required"`
		PageSize int `json:"pageSize" validate:"required"`
	}{}

	if err := context.BindAndValid(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	history, total, err := service.GetWeeklyEmailSendingHistory(req.Page, req.PageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, params.Response{
			Success: false,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, params.Response{
		Success: true,
		Data: map[string]interface{}{
			"total":   total,
			"history": history,
		},
	})
}
