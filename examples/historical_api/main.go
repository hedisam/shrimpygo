package main

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo"
	"github.com/hedisam/shrimpygo/examples/appconfig"
	"log"
	"time"
)

func main() {
	ctx := context.Background()
	cli := NewClient()
	historicalTrades(ctx, cli)
}

func historicalTrades(ctx context.Context, cli *shrimpygo.Client) {
	end := time.Now()
	start := end.Add(-1 * time.Hour * 24 * 7)
	trades, err := cli.GetHistoricalTrades(ctx, "binance", "btc-usdt", start, end, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Historical trades:", trades)
}

func historicalOrderBooks(ctx context.Context, cli *shrimpygo.Client) {
	end := time.Now()
	start := end.Add(-1 * time.Hour * 24 * 7)
	obs, err := cli.GetHistoricalOrderBooks(ctx, "coinbasepro", "btc-usd", start, end, 10)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Historical OrderBooks:", obs)
}

func historicalInstruments(ctx context.Context, cli *shrimpygo.Client) {
	inst, err := cli.GetHistoricalInstruments(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Historical Instruments:", inst)

	// there are three optional query parameters that you can provide: you can use the
	// utility function shrimpygo.QueryParams to build the query parameters
	var params []string

	// returning those instruments that are only available on a specific exchange
	// 'exchange' is the query parameter's name and 'coinbasepro' is the query parameter's value
	params = append(params, shrimpygo.QueryParams("exchange", "coinbasepro"))

	// returning only those instruments that match 'btc' as their base symbol
	params = append(params, shrimpygo.QueryParams("baseTradingSymbol", "btc"))

	// returning only those instruments that match 'usd' as their quote symbol; and maybe you want to use the
	// constant (shrimpygo.QuoteTradingSymbol or shrimpygo.BaseTradingSymbol) defined in the shrimpygo package as
	// query parameter's name to avoid misspellings.
	params = append(params, shrimpygo.QueryParams(shrimpygo.QuoteTradingSymbol, "usd"))

	// so now based on the query parameters that we've provided, the returned list of instruments should be limited to
	// the single pair 'btc-usd' from 'coinbasepro'

	inst, err = cli.GetHistoricalInstruments(ctx, params...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Instruments with query params:", inst)
}

func historicalCount(ctx context.Context, cli *shrimpygo.Client) {
	end := time.Now()
	start := end.Add(-1 * time.Hour * 24 * 7)
	dataType := "orderbook" // either 'orderbook' or 'trade'
	count, err := cli.GetHistoricalCount(ctx, "coinbasepro", dataType, "btc-usd", start, end)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Historical count: %d\n", count)
}

func historicalCandles(ctx context.Context, cli *shrimpygo.Client) {
	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Hour * 24 * 7)
	candles, err := cli.GetHistoricalCandles(ctx, "coinbasepro", "btc-usd", "1h", 100, startTime, endTime)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("historical candles:", candles)
}

func NewClient() *shrimpygo.Client {
	cfg, err := appconfig.Read("examples/appconfig/config.json")
	if err != nil {
		log.Fatal(err)
	}

	shrimpyCfg := shrimpygo.Config{PublicKey: cfg.APIKey, PrivateKey: cfg.SecretKey}
	client, err := shrimpygo.NewClient(shrimpyCfg)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
