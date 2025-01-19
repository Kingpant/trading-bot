package shared

import (
	"trading-bot/internal/config"

	log "github.com/sirupsen/logrus"
)

func NewLogger(AppEnv config.AppEnvironment) *log.Logger {
	logger := log.New()

	if AppEnv == config.Production {
		logger.SetFormatter(&log.JSONFormatter{})
		logger.SetLevel(log.InfoLevel)
	} else {
		logger.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	return logger
}
