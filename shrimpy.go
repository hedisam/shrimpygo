package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/ws"
	"strings"
	"time"
)

type Client struct {
	config *Config
}

func NewClient(cfx Config) (*Client, error) {
	if cfx.PublicApiKey() == "" || cfx.PrivateApiKey() == "" {
		return nil, fmt.Errorf("shrimpy config: invalid api/secret key")
	}

	return &Client{config: &cfx}, nil
}

// Websocket creates a websocket connection and returns a shrimpy WSConnection to interact with.
// throughput is just the downstream channel's capacity. Although in normal circumstances a small throughput works but
// since the websocket connection replaces unread messages whenever new data comes in, it's important to set an
// appropriate value as throughput to not miss any data. This is more important when you've subscribed to multiple
// channels of sequential data like order-book, especially in the times that market players are hyperactive.
func (cli *Client) Websocket(ctx context.Context, throughput int) (*WSConnection, error) {
	// connect to the ws server and create a ws stream
	stream, err := ws.CreateStream(ctx, throughput, cli.config)
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to create websocket stream: %w", err)
	}

	return &WSConnection{stream: stream, ctx: ctx, throughput: throughput}, nil
}

// SupportedExchanges (@PublicAPI) retrieves all Shrimpy supported exchanges and some basic information about each.
// no API keys are attached to the http request if you set freeApiCall which makes it a cost free
// request (rate limiting may be applied, check the docs @ https://developers.shrimpy.io/docs/#public)
func (cli *Client) SupportedExchanges(ctx context.Context, freeApiCall bool) ([]ExchangeInfo, error) {
	return supportedExchanges(cli, ctx, freeApiCall)
}

// ExchangeAssets (@PublicAPI) retrieves exchange asset information for a particular exchange.
// no API keys are attached to the http request if you set freeApiCall which makes it a cost free
// request (rate limiting may be applied, check the docs @ https://developers.shrimpy.io/docs/#public)
// Note that Shrimpy hosts a logo for each asset according to the table below:
// 32x32 	-> https://assets.shrimpy.io/cryptoicons/png/<id>.png
// 128x128 	-> https://assets.shrimpy.io/cryptoicons/png128/<id>.png
func (cli *Client) ExchangeAssets(ctx context.Context, exchange string, freeApiCall bool) ([]ExchangeAsset, error) {
	return exchangeAssets(cli, ctx, exchange, freeApiCall)
}

// TradingPairs (@PublicAPI) retrieves a list of active trading pairs for a particular exchange. The symbols will match the
// tradingSymbol from Get Exchange Assets as well as the symbol used by the exchange.
// no API keys are attached to the http request if you set freeApiCall which makes it a cost free
// request (rate limiting may be applied, check the docs @ https://developers.shrimpy.io/docs/#public)
func (cli *Client) TradingPairs(ctx context.Context, exchange string, freeApiCall bool) ([]TradingPair, error) {
	return tradingPairs(cli, ctx, exchange, freeApiCall)
}

// GetTicker retrieves all Shrimpy supported exchange assets for a particular exchange along with pricing information.
// Note: The symbol for the same asset may vary based on the exchange. For example, Stellar is "STR" on Poloniex,
// but "XLM" on other exchanges.
func (cli *Client) GetTicker(ctx context.Context, exchange string) ([]Ticker, error) {
	return getTicker(cli, ctx, exchange)
}

// GetOrderBooks retrieves live order book data. Examples of data that can be retrieved in a single call are below:
// * full depth order book for a single currency pair
// * order book for all currency pairs on a single exchange
// * order books for currency pairs across multiple exchanges
// See the MarketOrderBooks type for more information.
// exchanges:
//		The exchange for which to retrieve live order book data(e.g. coinbasepro);
// 		OR A comma separated list of exchanges (e.g. "bittrex,binance")
//		OR "all" to retrieve data for all supported exchanges
// queryOptions:
//		It could be a combining of these three: baseSymbol, quoteSymbol, or limit. They're all optional.
//		You can use the helper function QueryParams to build your query parameters. (see the exmaples)
// 		baseSymbol:
//			The base symbol. (e.g. XLM for a XLM-BTC market) quantity is in these units.
//			OR A comma separated list of baseSymbols (e.g. "XLM,STR,ETH")
//			Note: if baseSymbol is not supplied, all markets matching the quoteSymbol will be returned.
//		quoteSymbol:
// 			The quote symbol. (e.g. BTC for a XLM-BTC market)
//			OR A comma separated list of quoteSymbols (e.g. "BTC,USDT")
//			Note: if quoteSymbol is not supplied, all markets matching the baseSymbol will be returned.
//		limit:
//			The maximum number of asks and bids to retrieve. Defaults to 10 if not supplied.
//			Note: if requesting more than one market, limit cannot be greater than 10.
func (cli *Client) GetOrderBooks(ctx context.Context, exchanges string, queryOptions ...string) ([]MarketOrderBooks, error) {
	return getOrderBooks(cli, ctx, exchanges, queryOptions...)
}

// GetCandles retrieves live candlestick data. The candlestick data is typically used to plot a candlestick or OHLCV chart.
// When retrieving candlestick data for plotting, first call the endpoint without specifying a startTime. This will
// return data associated with the most recent 1000 candlesticks. Subsequently, periodically call the endpoint
// specifying the startTime as the time associated with the most recent candlestick. Note that the last or most recent
// candlestick is for the current, not-yet-committed frame.
// All times are in returned in UTC.
// Note: if no trades occur within a particular time frame, no candlestick data will be created for that time frame.
// pair: pair is the currency pair you're interested in, it could be either like "base/quote" or "base-quote" (e.g. BTC/USD)
// interval: one of these values: 1m, 5m, 15m, 1h, 6h, or 1d
// startTime(optional): Optionally only return data on or after the supplied startTime (inclusive).
func (cli *Client) GetCandles(ctx context.Context, exchange, pair, interval string,
	startTime ...string) ([]CandleStick, error) {

	var st string
	if len(startTime) == 1 {
		st = startTime[0]
	} else if len(startTime) > 1 {
		return nil, fmt.Errorf("shrimpygo failed to retrieve candlesticks list: too many options provided for startTime: please " +
			"provide max one option")
	}

	base, quote, err := pairSeparator(pair)
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve candlesticks list: %w", err)
	}

	return getCandles(cli, ctx, exchange, base, quote, interval, st)
}

// GetHistCandles retrieves historical candlestick (OHLCV) data based on the exchange, start and end times, and currency
// pair passed in. Data is returned as lists of HistoricalCandlestick objects.
// startTime and endTime need to be in ISO 8601 format.
// pair: the currency pair you're interested in, it could be either like "base/quote" or "base-quote" (e.g. BTC/USD)
// limit: the amount of items to return. Must be an integer from 1 to 1000.
// interval: must be one of the following values: 1m, 5m, 15m, 1h, 6h, or 1d
func (cli *Client) GetHistoricalCandles(ctx context.Context, exchange, pair string, interval string, limit int,
	startTime, endTime time.Time) ([]HistoricalCandlestick, error) {
	base, quote, err := pairSeparator(pair)
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical candlesticks list: %w", err)
	}
	return getHistoricalCandles(cli, ctx, exchange, base, quote, interval, limit, startTime, endTime)
}

// GetHistoricalCount retrieves that number of data points available based on the type, exchange, start and end times,
// and currency pair passed in. for dataType (type) of 'orderbook' the number of snapshots will be returned.
// startTime and endTime time must be rounded to the nearest hour (i.e. minutes and seconds must be 0)
// dataType: corresponds to the 'type' query parameter, either 'trade' or 'orderbook'.
func (cli *Client) GetHistoricalCount(ctx context.Context, exchange, dataType, pair string,
	startTime, endTime time.Time) (int, error) {
	base, quote, err := pairSeparator(pair)
	if err != nil {
		return -1, fmt.Errorf("shrimpygo failed to retrieve historical count: %w", err)
	}

	return getHistoricalCount(cli, ctx, exchange, dataType, base, quote, startTime, endTime)
}

// GetHistoricalInstruments retrieves the supported historical instruments on the Developer API and their available
// historical order book and trade data ranges.
// queryOptions:
//		It could be a combining of these three: baseTradingSymbol, quoteTradingSymbol, and/or exchange. They're all
//		optional.
//		You can use the helper function QueryParams to build your query parameters. (see the exmaples)
// 		exchange:
//			If provided, only trading pairs that are available on exchange will be returned
//		baseTradingSymbol:
//			If provided, only trading pairs with a matching base trading symbol will be returned
//		quoteTradingSymbol:
//			If provided, only trading pairs with a matching quote trading symbol will be returned
func (cli *Client) GetHistoricalInstruments(ctx context.Context, queryOptions ...string) ([]HistoricalInstrument, error) {
	return getHistoricalInstruments(cli, ctx, queryOptions...)
}

// GetHistoricalOrderBooks retrieves historical orderbooks based on the exchange, start and end times, and currency pair
// passed in. Data is returned as lists of HistoricalOrderBook objects as demonstrated in the examples.
// Historical Order Book Snapshots are taken roughly every minute.
// pair: the currency pair you're interested in, it could be either like "base/quote" or "base-quote" (e.g. BTC/USD)
// limit: the amount of items to return. Must be an integer from 1 to 1000.
func (cli *Client) GetHistoricalOrderBooks(ctx context.Context, exchange, pair string, startTime, endTime time.Time,
	limit int) ([]HistoricalOrderBook, error) {
	base, quote, err := pairSeparator(pair)
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical orderbooks: %w", err)
	}
	return getHistoricalOrderBooks(cli, ctx, exchange, base, quote, startTime, endTime, limit)
}

// GetHistoricalTrades retrieves historical trades based on the exchange, start and end times, and currency pair
// passed in. Data is returned as lists of HistoricalTrade objects as demonstrated in the examples.
// pair: the currency pair you're interested in, it could be either like "base/quote" or "base-quote" (e.g. BTC/USD)
// limit: the amount of items to return. Must be an integer from 1 to 1000.
func (cli *Client) GetHistoricalTrades(ctx context.Context, exchange, pair string, startTime, endTime time.Time,
	limit int) ([]HistoricalTrade, error) {
	base, quote, err := pairSeparator(pair)
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve historical trades: %w", err)
	}
	return getHistoricalTrades(cli, ctx, exchange, base, quote, startTime, endTime, limit)
}

func pairSeparator(pair string) (base, quote string, err error) {
	if strings.Contains("/", pair) {
		p := strings.Split(pair, "/")
		base = p[0]
		quote = p[1]
		return
	} else if strings.Contains(pair, "-") {
		p := strings.Split(pair, "-")
		base = p[0]
		quote = p[1]
		return
	}

	err = fmt.Errorf("failed to extract baseSybmol & quoteSymbol from the pair: 'pair' must be provided either like " +
		"'base/quote' or 'base-quote'")
	return
}
