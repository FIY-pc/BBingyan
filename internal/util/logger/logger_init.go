package logger

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
)

type logger struct {
	*zap.Logger
	traceId string
	spanId  string
	pSpanId string
}

var Logger logger

func InitLogger() {
	// 基础设置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	// 创建writer
	fileWriter := getFileLogWriter()
	cmdWriter := zapcore.AddSync(os.Stdout)

	var cores []zapcore.Core
	// 根据环境选择level
	switch config.Config.Server.Env {
	case "dev", "test":
		cores = append(cores, zapcore.NewCore(encoder, cmdWriter, zapcore.DebugLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, zapcore.DebugLevel))
	case "prod":
		cores = append(cores, zapcore.NewCore(encoder, cmdWriter, zapcore.InfoLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, zapcore.WarnLevel))
	default:
		cores = append(cores, zapcore.NewCore(encoder, cmdWriter, zapcore.DebugLevel))
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, zapcore.DebugLevel))
	}
	core := zapcore.NewTee(cores...)
	Logger = logger{Logger: zap.New(core)}
}

func getFileLogWriter() (writeSyncer zapcore.WriteSyncer) {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.Config.Log.LogFile,
		MaxSize:   config.Config.Log.MaxSize,
		MaxAge:    config.Config.Log.MaxAge,
		Compress:  config.Config.Log.Compress,
		LocalTime: config.Config.Log.LocalTime,
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
