package shrimpygo

type ExchangeInfo struct {
	Exchange     string  `json:"exchange"`
	BestCaseFee  float64 `json:"bestCaseFee"`
	WorstCaseFee float64 `json:"worstCaseFee"`
	IconUrl      string  `json:"icon"`
}

type ExchangeAsset struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	TradingSymbol string `json:"tradingSymbol"`
}

type TradingPair struct {
	BaseSymbol  string `json:"baseTradingSymbol"`
	QuoteSymbol string `json:"quoteTradingSymbol"`
}

/////////////////////////////////////////

type Ticker struct {
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	PriceUSD      string `json:"priceUsd"`
	PriceBTC      string `json:"priceBtc"`
	Last24hChange string `json:"percentChange24hUsd"`
	LastUpdated   string `json:"lastUpdated"`
}

type MarketOrderBooks struct {
	BaseSymbol  string              `json:"baseSymbol"`
	QuoteSymbol string              `json:"quoteSymbol"`
	OrderBooks  []ExchangeOrderBook `json:"orderBooks"`
}

type ExchangeOrderBook struct {
	Exchange  string    `json:"exchange"`
	OrderBook OrderBook `json:"orderBook"`
}

type OrderBook struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
}

type CandleStick struct {
	Open        string  `json:"open"`
	High        string  `json:"high"`
	Low         string  `json:"low"`
	Close       string  `json:"close"`
	Volume      string  `json:"volume"`
	QuoteVolume float64 `json:"quoteVolume"`
	BTCVolume   float64 `json:"btcVolume"`
	USDVolume   float64 `json:"usdVolume"`
	Time        string  `json:"time"`
}

/////////////////////////////////////////

type HistoricalCandlestick struct {
	Open        string  `json:"open"`
	High        string  `json:"high"`
	Low         string  `json:"low"`
	Close       string  `json:"close"`
	Volume      string  `json:"volume"`
	QuoteVolume float64 `json:"quoteVolume"`
	BTCVolume   float64 `json:"btcVolume"`
	USDVolume   float64 `json:"usdVolume"`
	Time        string  `json:"time"`
}

type HistoricalInstrument struct {
	Exchange    string `json:"exchange"`
	BaseSymbol  string `json:"baseTradingSymbol"`
	QuoteSymbol string `json:"quoteTradingSymbol"`
	// The start time of first datapoint for this instrument
	OrderBookStartTime string `json:"orderBookStartTime"`
	// The end time of last datapoint for this instrument
	OrderBookEndTime string `json:"orderBookEndTime"`
	// The start time of first datapoint for this instrument
	TradeStartTime string `json:"tradeStartTime"`
	// The end time of last datapoint for this instrument
	TradeEndTime string `json:"tradeEndTime"`
}

type HistoricalTrade struct {
	Time  string `json:"time"`
	Size  string `json:"size"`
	Price string `json:"price"`
	// Can be either 'buyer', 'seller' or 'unknown'.
	TakerSide string `json:"takerSide"`
}

type HistoricalOrderBook struct {
	// The id of the limit order !!! (from the docs: https://developers.shrimpy.io/docs/#historical-trade)
	Time string `json:"time"`
	// The best asks for the market. (ascending order)
	Asks []HistOrderBookItem `json:"asks"`
	// The best bids for the market. (descending order)
	Bids []HistOrderBookItem `json:"bids"`
}

type HistOrderBookItem struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

/////////////////////////////////////////

// OrderBookInfo returned by the websocket api
type OrderBookInfo struct {
	Exchange string           `json:"exchange"`
	Pair     string           `json:"pair"`
	Channel  string           `json:"channel"`
	Snapshot bool             `json:"snapshot"`
	Sequence int64            `json:"sequence"`
	Content  OrderBookContent `json:"content"`
}

type OrderBookContent struct {
	Sequence int64           `json:"sequence"`
	Asks     []OrderBookItem `json:"asks"`
	Bids     []OrderBookItem `json:"bids"`
}

type OrderBookItem struct {
	Price    string `json:"Price"`
	Quantity string `json:"quantity"`
}

/////////////////////////////////////////

type Trades struct {
	Exchange string       `json:"exchange"`
	Pair     string       `json:"pair"`
	Channel  string       `json:"channel"`
	Snapshot bool         `json:"snapshot"`
	Sequence int64        `json:"sequence"`
	Content  []TradesItem `json:"content"`
}

type TradesItem struct {
	Id        string  `json:"id"`
	Price     string  `json:"price"`
	Quantity  string  `json:"quantity"`
	Time      string  `json:"time"`
	BTCValue  float64 `json:"btcValue"`
	USDValue  float64 `json:"usdValue"`
	TakerSide string  `json:"takerSide"`
}

////////////////////////////////////////

type Orders struct {
	Channel string   `json:"channel"`
	Content []string `json:"content"`
}

///////////////////////////////////////

type Subscription struct {
	Type     string `json:"type"`
	Exchange string `json:"exchange"`
	Pair     string `json:"pair"`
	Channel  string `json:"channel"`
}

type pingPong struct {
	Type string `json:"type"`
	Data int64  `json:"data"`
}

type Error struct {
	Type    string `json:"type"`
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// unknownData streamed by the ws connection
type unknownData struct {
	// having n Channel field shows it's not a ping nor an error. then we can decode the data based on the channel.
	Channel string `json:"channel"`
	// both error and ping messages come with a Type field
	Type string `json:"type"`
	// ping messages come with a Data field
	Data int64 `json:"data"`
	// Code and Message are expected with an error message
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
