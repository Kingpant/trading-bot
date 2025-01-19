package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"trading-bot/internal/config"
	"trading-bot/internal/connectors/bitkub/api"
	"trading-bot/internal/modules/bitkub"
	"trading-bot/internal/shared"

	"github.com/sirupsen/logrus"
)

func main() {
	_, ctxCancel := context.WithCancel(context.Background())

	cfg, loadConfigErr := config.LoadConfig(os.Getenv("DOTENV_PATH"))
	if loadConfigErr != nil {
		log.Panicln(loadConfigErr)
	}

	logger := shared.NewLogger(cfg.AppEnv)

	bitkubAPIClient := api.NewBitkubApiClient(cfg.BitkubBaseURL, cfg.BitkubAPIKey, cfg.BitkubAPISecret)

	_ = bitkub.NewBitkubModule(bitkubAPIClient)

	gracefulShutdown(ctxCancel, logger)

	logger.Info("App started")
}

func gracefulShutdown(ctxCancel context.CancelFunc, logger *logrus.Logger) {
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop
	logger.Info("App is shutting down...")
	ctxCancel()
}
