package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"trading-bot/internal/connectors/bitkub/api/dto"
)

type IBitkubApiClient interface {
	RequestDepositHistories(ctx context.Context) ([]dto.DepositHistory, error)
	RequestOrderHistories(
		ctx context.Context,
		tokenSymbol string,
		startTimestamp *uint64,
	) ([]dto.OrderHistory, error)
}

type bitkubApiClient struct {
	baseUrl   string
	apiKey    string
	apiSecret string
}

func NewBitkubApiClient(baseUrl, apiKey, apiSecret string) IBitkubApiClient {
	return &bitkubApiClient{
		baseUrl:   baseUrl,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (b *bitkubApiClient) RequestDepositHistories(
	ctx context.Context,
) ([]dto.DepositHistory, error) {
	page := "1"
	limit := "100"
	depositHistories := []dto.DepositHistory{}

	for {
		depositHistoryPath := "/api/v3/crypto/deposit-history"
		queryParma := fmt.Sprintf("?p=%s&lmt=%s", page, limit)

		depositHistoryResponse := dto.DepositHistoryResponse{}
		err := b.httpGetRequest(depositHistoryPath, queryParma, &depositHistoryResponse)
		if err != nil {
			return nil, err
		}

		depositHistories = append(depositHistories, depositHistoryResponse.Result...)

		if depositHistoryResponse.Pagination.Next == 0 {
			break
		}

		page = strconv.FormatUint(depositHistoryResponse.Pagination.Next, 10)
	}

	return depositHistories, nil
}

func (b *bitkubApiClient) RequestOrderHistories(
	ctx context.Context,
	tokenSymbol string,
	startTimestamp *uint64,
) ([]dto.OrderHistory, error) {
	tokenSymbol = strings.ToLower(tokenSymbol)
	page := "1"
	limit := "100"
	orderHistories := []dto.OrderHistory{}

	for {
		orderHistoryPath := "/api/v3/market/my-order-history"
		queryParma := fmt.Sprintf("?sym=%s_thb&p=%s&lmt=%s", tokenSymbol, page, limit)
		if startTimestamp != nil {
			queryParma += fmt.Sprintf("&start=%d", *startTimestamp)
		}

		orderHistoryResponse := dto.OrderHistoryResponse{}
		err := b.httpPostRequest(orderHistoryPath, queryParma, &orderHistoryResponse)
		if err != nil {
			return nil, err
		}

		orderHistories = append(orderHistories, orderHistoryResponse.Result...)

		if orderHistoryResponse.Pagination.Next == 0 {
			break
		}
		page = strconv.FormatUint(orderHistoryResponse.Pagination.Next, 10)
	}

	return orderHistories, nil
}

func (b *bitkubApiClient) httpGetRequest(path, queryParam string, respType interface{}) error {
	nowMilliSec := time.Now().UnixMilli()
	nowMilliSecStr := strconv.FormatInt(nowMilliSec, 10)
	method := "GET"
	payload := nowMilliSecStr + method + path + queryParam

	signature := b.genSignature(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, b.baseUrl+path+queryParam, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-BTK-TIMESTAMP", nowMilliSecStr)
	req.Header.Add("X-BTK-APIKEY", b.apiKey)
	req.Header.Add("X-BTK-SIGN", signature)

	resp, err := client.Do(req)
	if err != nil {
		return err
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

func (b *bitkubApiClient) httpPostRequest(
	path string,
	queryParam string,
	respType interface{},
) error {
	nowMilliSec := time.Now().UnixMilli()
	nowMilliSecStr := strconv.FormatInt(nowMilliSec, 10)
	method := "POST"
	payload := nowMilliSecStr + method + path + queryParam

	signature := b.genSignature(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, b.baseUrl+path+queryParam, nil)
	if err != nil {
		return err
	}

	req.Header.Add("X-BTK-TIMESTAMP", nowMilliSecStr)
	req.Header.Add("X-BTK-APIKEY", b.apiKey)
	req.Header.Add("X-BTK-SIGN", signature)

	resp, err := client.Do(req)
	if err != nil {
		return err
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

func (b *bitkubApiClient) genSignature(payload string) string {
	secretBytes := []byte(b.apiSecret)
	payloadBytes := []byte(payload)

	hmac := hmac.New(sha256.New, secretBytes)
	hmac.Write(payloadBytes)
	signature := hmac.Sum(nil)

	return hex.EncodeToString(signature)
}
