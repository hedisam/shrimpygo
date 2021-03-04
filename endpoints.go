package shrimpygo

const (
	supportedExchangesApi = "/v1/list_exchanges"
	exchangeAssetsApi     = "/v1/exchanges/%s/assets"        // %s to be replaced by the exchange id
	tradingPairsApi       = "/v1/exchanges/%s/trading_pairs" // %s to be replaced by the exchange id
)
