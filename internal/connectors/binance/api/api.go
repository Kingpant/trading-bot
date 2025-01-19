package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"trading-bot/internal/connectors/binance/api/dto"
)

type IBinanceAPI interface {
	Ping(ctx context.Context) error
	CheckServerTime(ctx context.Context) (*dto.CheckServerTimeResponse, error)
}

type binanceAPI struct {
	BaseURL string
}

func NewBinanceAPI(baseURL string) IBinanceAPI {
	return &binanceAPI{
		BaseURL: baseURL,
	}
}

func (b *binanceAPI) Ping(ctx context.Context) error {
	pingErr := b.sendGetRequest(ctx, "/fapi/v1/ping", nil)
	return pingErr
}

func (b *binanceAPI) CheckServerTime(ctx context.Context) (*dto.CheckServerTimeResponse, error) {
	var resp dto.CheckServerTimeResponse
	checkServerTimeErr := b.sendGetRequest(ctx, "/fapi/v1/time", &resp)
	return &resp, checkServerTimeErr
}

func (b *binanceAPI) sendGetRequest(ctx context.Context, path string, respType interface{}) error {
	req, newRequestErr := http.NewRequestWithContext(ctx, http.MethodGet, b.BaseURL+path, nil)
	if newRequestErr != nil {
		return newRequestErr
	}

	resp, doErr := http.DefaultClient.Do(req)
	if doErr != nil {
		return doErr
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(respType)
	if decodeErr != nil {
		return decodeErr
	}

	defer resp.Body.Close()

	return nil
}
