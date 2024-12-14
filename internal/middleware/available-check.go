package middleware

import (
	"github.com/FIY-pc/BBingyan/internal/controller/params"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/es"
	"github.com/FIY-pc/BBingyan/internal/infrastructure/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

// EsAvailableCheck checks if the ES service is available
func EsAvailableCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if es.ES == nil {
			logger.Log.Error(c, "ES service is not available")
			return c.JSON(http.StatusBadGateway, params.Response{
				Success: false,
				Message: "ES service is not available",
			})
		}
		return next(c)
	}
}
