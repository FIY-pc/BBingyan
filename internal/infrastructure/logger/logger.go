package logger

import (
	"fmt"
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
)

var Log *Logger

type Logger struct {
	*zap.Logger
	traceId string
	spanId  string
	pSpanId string
}

func NewLogger() {
	// 基础设置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	// 创建writer
	fileWriter := getFileLogWriter()
	cmdWriter := zapcore.AddSync(os.Stdout)

	var cores []zapcore.Core
	// 根据环境选择level
	switch config.Configs.Server.Env {
	case "dev", "test":
		cores = append(cores, zapcore.NewCore(encoder, cmdWriter, zapcore.DebugLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, zapcore.DebugLevel))
	case "prod":
		cores = append(cores, zapcore.NewCore(encoder, cmdWriter, zapcore.WarnLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, zapcore.ErrorLevel))
	default:
		cores = append(cores, zapcore.NewCore(encoder, cmdWriter, zapcore.DebugLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, zapcore.DebugLevel))
	}
	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core)
	// 注入Log全局变量
	Log = &Logger{Logger: zapLogger}
}

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.Configs.Log.LogFile,
		MaxSize:   config.Configs.Log.MaxSize,
		MaxAge:    config.Configs.Log.MaxAge,
		Compress:  config.Configs.Log.Compress,
		LocalTime: config.Configs.Log.LocalTime,
	}
	writeSyncer = zapcore.AddSync(lumberJackLogger)
	return writeSyncer
}

func GetLogCallerInfo() (file, funcName string, line int) {
	// 获取第三层调用栈的信息
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return
	}
	// 获取文件名
	file = path.Base(file)
	// 获取调用函数名
	funcName = runtime.FuncForPC(pc).Name()
	return file, funcName, line
}

func (l *Logger) Debug(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.DebugLevel, msg, kv)
}

func (l *Logger) Info(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.InfoLevel, msg, kv)
}

func (l *Logger) Warn(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.WarnLevel, msg, kv)
}

func (l *Logger) Error(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.ErrorLevel, msg, kv)
}

func (l *Logger) Panic(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.PanicLevel, msg, kv)
}

func (l *Logger) Fatal(c echo.Context, msg string, kv ...interface{}) {
	if l == nil {
		return
	}
	l.getTrace(c)
	l.log(zapcore.FatalLevel, msg, kv)
}

func (l *Logger) log(lvl zapcore.Level, msg string, kv ...interface{}) {
	if len(kv)%2 != 0 {
		kv = append(kv, "unknown")
	}
	// 检查日志级别
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
	if ce != nil {
		ce.Write(Fields...)
	}
}

func (l *Logger) getTrace(c echo.Context) {
	if l == nil {
		return
	}
	if c == nil {
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
		l.pSpanId = pSpanId
	}
}
