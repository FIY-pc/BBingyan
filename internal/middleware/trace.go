package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TraceMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		traceId := c.Request().Header.Get("X-Trace-ID")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		spanId := uuid.New().String()[:8]
		pSpanId := c.Request().Header.Get("X-Span-ID")

		c.Request().Header.Set("X-Trace-ID", traceId)
		c.Request().Header.Set("X-Span-ID", spanId)
		c.Request().Header.Set("X-PSpan-ID", pSpanId)
		return next(c)
	}
}
