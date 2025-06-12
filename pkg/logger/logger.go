package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)
	return logger
}

const LoggerContextKey string = "logger"
