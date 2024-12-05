package logger

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (l *logger) Debug(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.InfoLevel, msg, kv)
}

func (l *logger) Info(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.InfoLevel, msg, kv)
}

func (l *logger) Warn(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.WarnLevel, msg, kv)
}

func (l *logger) Error(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.ErrorLevel, msg, kv)
}

func (l *logger) Panic(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.PanicLevel, msg, kv)
	panic(msg)
}

func (l *logger) log(lvl zapcore.Level, msg string, kv ...interface{}) {
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	ce := l.Check(lvl, msg)
	// 加入接口调用信息
	kv = append(kv, "traceId", l.traceId, "spanId", l.spanId, "pSpanId", l.pSpanId)
	// 加入调用栈信息
	file, funcName, line := GetLogCallerInfo()
	kv = append(kv, "file", file, "funcName", funcName, "line", line)
	// 将所有信息读入Field中
	Fields := make([]zap.Field, 0, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k := fmt.Sprintf("%v", kv[i])
		Fields = append(Fields, zap.Any(k, kv[i+1]))
	}
	ce.Write(Fields...)
}

func (l *logger) getTrace(c echo.Context) {
	if l == nil {
		return
	}
	traceId := c.Request().Header.Get("X-Trace-ID")
	spanId := c.Request().Header.Get("X-Span-ID")
	pSpanId := c.Request().Header.Get("X-PSpan-ID")
	if traceId != "" {
		l.traceId = traceId
	}
	if spanId != "" {
		l.spanId = spanId
	}
	if pSpanId != "" {
		l.pSpanId = spanId
	}
}
