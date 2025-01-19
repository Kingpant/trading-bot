package binance

import "context"

func (b *binanceModule) ping(ctx context.Context) error {
	return b.binanceAPIClient.Ping(ctx)
}

func (b *binanceModule) checkServerTime(ctx context.Context) error {
	_, checkServerTimeErr := b.binanceAPIClient.CheckServerTime(ctx)
	return checkServerTimeErr
}
