package bitkub

import (
	"trading-bot/internal/connectors/bitkub/api"

	"github.com/sirupsen/logrus"
)

type IBitkubModule interface {
}

type bitkubModule struct {
	bitkubAPIClient api.IBitkubApiClient
	logger          *logrus.Logger
}

func NewBitkubModule(bitkubAPIClient api.IBitkubApiClient, logger *logrus.Logger) IBitkubModule {
	return &bitkubModule{
		bitkubAPIClient: bitkubAPIClient,
		logger:          logger,
	}
}
