package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/rest"
)

func getTicker(cli *Client, ctx context.Context, exchange string) ([]Ticker, error) {
	var tickers []Ticker

	err := rest.HttpGet(ctx, fmt.Sprintf(getTickerApi, exchange), cli.config, rest.NewDecoderFunc(&tickers))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve tickers list: %w", err)
	}

	return tickers, nil
}

func getOrderBooks(cli *Client, ctx context.Context, exchanges string, queryOptions ...string) ([]MarketOrderBooks, error) {
	var orderBooks []MarketOrderBooks

	path := fmt.Sprint(getOrderBooksApi, "?exchange=", exchanges)
	for _, option := range queryOptions {
		path = fmt.Sprint(path, "&", option)
	}

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&orderBooks))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve market orderbooks list: %w", err)
	}

	return orderBooks, nil
}

func getCandles(cli *Client, ctx context.Context, exchange, baseSymbol, quoteSymbol, interval, startTime string) ([]CandleStick, error) {
	var candles []CandleStick
	path := fmt.Sprintf(getCandlesApi, exchange)
	path = fmt.Sprint(path, "?baseTradingSymbol=", baseSymbol, "&quoteTradingSymbol=", quoteSymbol, "&interval=", interval)
	if startTime != "" {
		path += "&startTime=" + startTime
	}

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&candles))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve candlesticks list: %w", err)
	}

	return candles, nil
}


