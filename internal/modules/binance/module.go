package binance

import (
	"trading-bot/internal/connectors/binance/api"
)

type binanceModule struct {
	binanceAPIClient api.IBinanceAPI
}

func NewBinanceModule(binanceAPIClient api.IBinanceAPI) *binanceModule {
	return &binanceModule{
		binanceAPIClient: binanceAPIClient,
	}
}
