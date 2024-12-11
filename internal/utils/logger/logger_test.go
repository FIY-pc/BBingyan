package logger

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	// 创建一个测试配置
	config.LoadConfig()

	// 初始化 Logger
	NewLogger()
	// 创建一个 Echo 上下文
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Trace-ID", "test-trace-id")
	req.Header.Set("X-Span-ID", "test-span-id")
	req.Header.Set("X-PSpan-ID", "test-pspan-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// 测试不同的日志级别
	tests := []struct {
		level   string
		message string
	}{
		{"Debug", "This is a debug message"},
		{"Info", "This is an info message"},
		{"Warn", "This is a warning message"},
		{"Error", "This is an error message"},
	}

	for _, tt := range tests {
		switch tt.level {
		case "Debug":
			Log.Debug(c, tt.message)
		case "Info":
			Log.Info(c, tt.message)
		case "Warn":
			Log.Warn(c, tt.message)
		case "Error":
			Log.Error(c, tt.message)
		}

	}
}
