package logger

import (
	"os"

	"github.com/sirupsen/logrus"

	"emarcey/data-vault/common"
)

func MakeLogger(loggerType string, env string) (*logrus.Logger, error) {
	severityLevel := logrus.DebugLevel
	if env != "local" {
		severityLevel = logrus.InfoLevel
	}
	logger := logrus.New()
	switch loggerType {
	case "json":
		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetOutput(os.Stdout)
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		return nil, common.NewInitializationError("logger", "Unknown logger type %s", loggerType)
	}
	logger.SetLevel(severityLevel)
	return logger, nil
}
