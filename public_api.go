package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/rest"
)

func supportedExchanges(cli *Client, ctx context.Context, freeApiCall bool) ([]ExchangeInfo, error) {
	var exchanges []ExchangeInfo

	err := publicApi(cli, ctx, supportedExchangesApi, freeApiCall, rest.NewDecoderFunc(&exchanges))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve the supported exchanges list: %w", err)
	}

	return exchanges, nil
}

func exchangeAssets(cli *Client, ctx context.Context, exchange string, freeApiCall bool) ([]ExchangeAsset, error) {
	var assets []ExchangeAsset

	err := publicApi(cli, ctx, fmt.Sprintf(exchangeAssetsApi, exchange), freeApiCall, rest.NewDecoderFunc(&assets))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve exchange assets list: %w", err)
	}

	return assets, nil
}

func tradingPairs(cli *Client, ctx context.Context, exchange string, freeApiCall bool) ([]TradingPair, error) {
	var pairs []TradingPair

	err := publicApi(cli, ctx, fmt.Sprintf(tradingPairsApi, exchange), freeApiCall, rest.NewDecoderFunc(&pairs))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve trading pairs list: %w", err)
	}

	return pairs, nil
}

func publicApi(cli *Client, ctx context.Context, urlPath string, freeApiCall bool, decoder rest.Decoder) error {
	var cfg rest.Config = cli.config
	if freeApiCall {
		cfg = nil
	}
	return rest.HttpGet(ctx, urlPath, cfg, decoder)
}
