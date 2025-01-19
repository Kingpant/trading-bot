package bitkub

import "trading-bot/internal/connectors/bitkub/api"

type bitkubModule struct {
	bitkubAPIClient api.IBitkubApiClient
}

func NewBitkubModule(bitkubAPIClient api.IBitkubApiClient) *bitkubModule {
	return &bitkubModule{
		bitkubAPIClient: bitkubAPIClient,
	}
}
