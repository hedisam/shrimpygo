package shrimpygo

type ExchangeInfo struct {
	Exchange string `json:"exchange"`
	BestCaseFee float64 `json:"bestCaseFee"`
	WorstCaseFee float64 `json:"worstCaseFee"`
	IconUrl string `json:"icon"`
}

/////////////////////////////////////////

type OrderBook struct {
	Exchange string           `json:"exchange"`
	Pair     string           `json:"pair"`
	Channel  string           `json:"channel"`
	Snapshot bool             `json:"snapshot"`
	Sequence int64            `json:"sequence"`
	Content  OrderBookContent `json:"content"`
}

type OrderBookContent struct {
	Asks []OrderBookItem `json:"asks"`
	Bids []OrderBookItem `json:"bids"`
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

type unknownData struct {
	// having an Channel field shows it's not a ping nor an error. then we can decode the data based on the channel.
	Channel string `json:"channel"`
	// both error and ping messages come with a Type field
	Type string `json:"type"`
	// ping messages come with a Data field
	Data int64 `json:"data"`
	// Code and Message are expected with an error message
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
