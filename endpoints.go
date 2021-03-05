package shrimpygo

const (
	// public api
	supportedExchangesApi = "/v1/list_exchanges"
	exchangeAssetsApi     = "/v1/exchanges/%s/assets"        // %s to be replaced by the exchange id
	tradingPairsApi       = "/v1/exchanges/%s/trading_pairs" // %s to be replaced by the exchange id

	// market data api
	getTickerApi = "/v1/exchanges/%s/ticker" // %s to be replaced by the exchange id
	getOrderBooksApi = "/v1/orderbooks"
    getCandlesApi = "/v1/exchanges/%s/candles" // %s to be replaced by the exchange id
)