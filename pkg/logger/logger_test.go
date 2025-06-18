package logger_test

import (
	"testing"

	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger_DefaultLevel(t *testing.T) {
	log := logger.NewLogger("undefined")

	assert.Equal(t, logrus.InfoLevel, log.GetLevel())
	assert.IsType(t, &logrus.JSONFormatter{}, log.Formatter)
}

func TestNewLogger_VariousLevels(t *testing.T) {
	levels := map[string]logrus.Level{
		"trace": logrus.TraceLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}

	for levelStr, expectedLevel := range levels {
		log := logger.NewLogger(levelStr)
		assert.Equal(t, expectedLevel, log.GetLevel())
	}
}
