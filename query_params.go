package shrimpygo

// query parameter names
const (
	BaseSymbol = "baseSymbol"
	QuoteSymbol = "quoteSymbol"
	Limit 	= "limit"
	Exchange = "exchange"
	BaseTradingSymbol = "baseTradingSymbol"
	QuoteTradingSymbol = "quoteTradingSymbol"
)

func QueryParams(name string, params ...string) string {
	query := name + "="
	for i, param := range params {
		query += param
		if i < len(params) - 1 {
			query += ","
		}
	}
	return query
}
