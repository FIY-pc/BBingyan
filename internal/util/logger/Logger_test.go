package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 创建测试用的 logger 实例
func setupTestLogger(t *testing.T) *logger {
	// 创建测试用的 zap logger
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	zapLogger, err := cfg.Build()
	assert.NoError(t, err)

	return &logger{
		Logger: zapLogger,
	}
}

// 创建测试用的 echo.Context
func setupEchoContext(headers map[string]string) echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return e.NewContext(req, rec)
}

func TestLogger_GetTrace(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    struct {
			traceId string
			spanId  string
			pSpanId string
		}
	}{
		{
			name: "with all headers",
			headers: map[string]string{
				"X-Trace-ID": "trace123",
				"X-Span-ID":  "span123",
				"X-PSpan-ID": "pspan123",
			},
			want: struct {
				traceId string
				spanId  string
				pSpanId string
			}{
				traceId: "trace123",
				spanId:  "span123",
				pSpanId: "span123",
			},
		},
		{
			name:    "without headers",
			headers: map[string]string{},
			want: struct {
				traceId string
				spanId  string
				pSpanId string
			}{
				traceId: "",
				spanId:  "",
				pSpanId: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := setupTestLogger(t)
			c := setupEchoContext(tt.headers)

			l.getTrace(c)

			assert.Equal(t, tt.want.traceId, l.traceId)
			assert.Equal(t, tt.want.spanId, l.spanId)
			assert.Equal(t, tt.want.pSpanId, l.pSpanId)
		})
	}
}

func TestLogger_Log_Levels(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(*logger, echo.Context, string, ...interface{})
		level   zapcore.Level
	}{
		{
			name:    "Debug level",
			logFunc: (*logger).Debug,
			level:   zapcore.DebugLevel,
		},
		{
			name:    "Info level",
			logFunc: (*logger).Info,
			level:   zapcore.InfoLevel,
		},
		{
			name:    "Warn level",
			logFunc: (*logger).Warn,
			level:   zapcore.WarnLevel,
		},
		{
			name:    "Error level",
			logFunc: (*logger).Error,
			level:   zapcore.ErrorLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := setupTestLogger(t)
			c := setupEchoContext(map[string]string{
				"X-Trace-ID": "trace123",
				"X-Span-ID":  "span123",
				"X-PSpan-ID": "pspan123",
			})

			// 不应该panic
			assert.NotPanics(t, func() {
				tt.logFunc(l, c, "test message", "key1", "value1")
			})
		})
	}
}

func TestLogger_Panic(t *testing.T) {
	l := setupTestLogger(t)
	c := setupEchoContext(map[string]string{})

	assert.Panics(t, func() {
		l.Panic(c, "panic message")
	})
}

func TestLogger_NilLogger(t *testing.T) {
	var l *logger
	c := setupEchoContext(map[string]string{})

	// 不应该panic
	assert.NotPanics(t, func() {
		l.Debug(c, "test")
		l.Info(c, "test")
		l.Warn(c, "test")
		l.Error(c, "test")
	})
}

func TestLogger_Log_OddKeyValues(t *testing.T) {
	l := setupTestLogger(t)
	c := setupEchoContext(map[string]string{})

	// 测试奇数个参数的情况
	assert.NotPanics(t, func() {
		l.Info(c, "test message", "key1", "value1", "key2")
	})
}

func TestGetLogCallerInfo(t *testing.T) {
	file, funcName, line := GetLogCallerInfo()

	assert.NotEmpty(t, file)
	assert.NotEmpty(t, funcName)
	assert.Greater(t, line, 0)
}
