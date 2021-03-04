package shrimpygo

const (
	supportedExchanges = "/v1/list_exchanges"
	exchangeAssets = "/v1/exchanges/%s/assets" // %s to be replaced by the exchange id
	tradingPairs = "/v1/exchanges/%s/trading_pairs" // %s to be replaced by the exchange id
)
