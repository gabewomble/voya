package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

var loggerInstance *Logger

func New() *Logger {
	if loggerInstance != nil {
		return loggerInstance
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	loggerInstance = &Logger{
		SugaredLogger: logger.Sugar(),
	}

	return loggerInstance
}

func (l *Logger) LogError(c *gin.Context, message string, err error, fields ...interface{}) {
	baseFields := []interface{}{
		"method", c.Request.Method,
		"uri", c.Request.RequestURI,
		"error", err.Error(),
	}
	allFields := append(baseFields, fields...)
	l.Errorw("Error", allFields...)
}

func (l *Logger) LogInfo(c *gin.Context, message string, fields ...interface{}) {
	baseFields := []interface{}{
		"method", c.Request.Method,
		"uri", c.Request.RequestURI,
	}
	allFields := append(baseFields, fields...)
	l.Infow(message, allFields...)
}
