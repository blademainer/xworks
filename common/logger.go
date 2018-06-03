package common

import (
	"os"
	logger "github.com/sirupsen/logrus"
)

const (
	ENV_LOG_LEVEL   = "LOG_LEVEL"

	DEFAULT_LOG_LEVEL = "debug"

	LOG_LEVEL_DEBUG = "debug"
	LOG_LEVEL_INFO  = "info"
	LOG_LEVEL_WARN  = "warn"
	LOG_LEVEL_ERROR = "error"
)

func SetLogLevel() {
	logLevel, _ := os.LookupEnv(ENV_LOG_LEVEL)
	switch logLevel {
	case LOG_LEVEL_DEBUG:
		logger.SetLevel(logger.DebugLevel)
	case LOG_LEVEL_INFO:
		logger.SetLevel(logger.InfoLevel)
	case LOG_LEVEL_WARN:
		logger.SetLevel(logger.WarnLevel)
	case LOG_LEVEL_ERROR:
		logger.SetLevel(logger.ErrorLevel)
	default:
		logger.SetLevel(logger.DebugLevel)
	}
}
