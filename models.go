package shrimpygo

type PriceQuote struct {
	Price string `json:"Price"`
	Quantity string `json:"quantity"`
}

type PriceContent struct {
	Asks []PriceQuote `json:"asks"`
	Bids []PriceQuote `json:"bids"`
}

type PriceData struct {
	Exchange string `json:"exchange"`
	Pair string `json:"pair"`
	Channel string `json:"channel"`
	Snapshot bool `json:"snapshot"`
	Sequence int64 `json:"sequence"`
	Content PriceContent `json:"content"`
}

type Subscription struct {
	Type string `json:"type"`
	Exchange string `json:"exchange"`
	Pair string `json:"pair"`
	Channel string `json:"channel"`
}

type pingPong struct {
	Type string `json:"type"`
	Data int64 `json:"data"`
}

type Error struct {
	Type string `json:"type"`
	Code int64 `json:"code"`
	Message string `json:"message"`
}

type unknownData struct {
	// having an Exchange field shows that we have a Price message
	Exchange string `json:"exchange"`
	// both error and ping messages come with a Type field
	Type string `json:"type"`
	// ping pongListener messages come with a Data field
	Data int64 `json:"data"`
	// Code and Message are expected with an error message
	Code int64 `json:"code"`
	Message string `json:"message"`
}

type wsToken struct {
	Token string `json:"token"`
}