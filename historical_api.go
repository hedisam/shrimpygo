package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/rest"
	"time"
)

/*
	The Historical Data endpoints are used to retrieve historical data such as order books and trades.
	Master API keys with the Data permission enabled can use Historical endpoints.
*/

type historicalCount struct {
	Count int `json:"count"`
}

func getHistoricalCandles(cli *Client, ctx context.Context, exchange, baseSymbol, quoteSymbol, interval string, limit int,
	startTime, endTime time.Time) ([]HistoricalCandlestick, error) {
	var candles []HistoricalCandlestick


	path := fmt.Sprintf("%s?exchange=%s&baseTradingSymbol=%s&quoteTradingSymbol=%s&startTime=%v&endTime=%v"+
		"&limit=%d&interval=%s", getHistCandlesApi, exchange, baseSymbol, quoteSymbol,
		timeToDate(startTime), timeToDate(endTime), limit, interval)

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&candles))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical candlesticks list: %w", err)
	}

	return candles, nil
}

func getHistoricalCount(cli *Client, ctx context.Context, exchange, dataType, baseSymbol, quoteSymbol string,
	startTime, endTime time.Time) (int, error) {
	var count historicalCount
	path := fmt.Sprintf("%s?type=%s&exchange=%s&baseTradingSymbol=%s&quoteTradingSymbol=%s&startTime=%v&endTime=%v",
		getHistCountApi, dataType, exchange, baseSymbol, quoteSymbol, timeToDate(startTime), timeToDate(endTime))

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&count))
	if err != nil {
		return -1, fmt.Errorf("shrimpygo failed to retrieve historical count: %w", err)
	}

	return count.Count, nil
}

func getHistoricalInstruments(cli *Client, ctx context.Context, queryOptions ...string) ([]HistoricalInstrument, error) {
	var instruments []HistoricalInstrument
	path := getHistInstApi

	for i, option := range queryOptions {
		if i == 0 {
			path += "?" + option
			continue
		}
		path += "&" + option
	}

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&instruments))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical instruments: %w", err)
	}

	return instruments, nil
}

func getHistoricalOrderBooks(cli *Client, ctx context.Context, exchange, baseSymbol, quoteSymbol string,
	startTime, endTime time.Time, limit int) ([]HistoricalOrderBook, error) {
	var orderBooks []HistoricalOrderBook
	path := fmt.Sprintf("%s?exchange=%s&baseTradingSymbol=%s&quoteTradingSymbol=%s&startTime=%v&endTime=%v&limit=%d",
		getHistOrderBooksApi, exchange, baseSymbol, quoteSymbol, timeToDate(startTime), timeToDate(endTime), limit)

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&orderBooks))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical orderbooks: %w", err)
	}

	return orderBooks, nil
}

func getHistoricalTrades(cli *Client, ctx context.Context, exchange, baseSymbol, quoteSymbol string,
	startTime, endTime time.Time, limit int) ([]HistoricalTrade, error) {
	var trades []HistoricalTrade
	path := fmt.Sprintf("%s?exchange=%s&baseTradingSymbol=%s&quoteTradingSymbol=%s&startTime=%v&endTime=%v&limit=%d",
		getHistTradesApi, exchange, baseSymbol, quoteSymbol, timeToDate(startTime), timeToDate(endTime), limit)

	err := rest.HttpGet(ctx, path, cli.config, rest.NewDecoderFunc(&trades))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical trades: %w", err)
	}

	return trades, nil
}

func timeToDate(t time.Time) string {
	y, m, d := t.Date()
	return fmt.Sprintf("%d-%d-%d", y, m, d)
}