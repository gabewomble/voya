package logger

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	*slog.Logger
}

var loggerInstance *Logger

func New() *Logger {
	if loggerInstance != nil {
		return loggerInstance
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	loggerInstance = &Logger{
		logger,
	}

	return loggerInstance
}

func (l *Logger) LogError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.RequestURI
	)

	l.Error(err.Error(), "method", method, "uri", uri)
}

func (l *Logger) LogInfo(c *gin.Context, message string) {
	var (
		method = c.Request.Method
		uri    = c.Request.RequestURI
	)
	l.Info(message, "method", method, "uri", uri)
}
