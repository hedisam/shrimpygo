package shrimpygo

const (
	ChannelBBO       = "bbo"
	ChannelOrderBook = "orderbook"
	ChannelTrades    = "trade"
	ChannelOrders    = "orders"
)

// BBOSubs returns a subscription of channel 'bbo'. Subscription's type will be set by Subscribe/Unsubscribe methods of
// WSConnection instance.
func BBOSubs(exchange, pair string) Subscription {
	return Subscription{
		Exchange: exchange,
		Pair:     pair,
		Channel:  ChannelBBO,
	}
}

// OrderBookSubs returns a subscription of channel 'orderbook'. Subscription's type will be set by Subscribe/Unsubscribe methods of
// WSConnection instance.
func OrderBookSubs(exchange, pair string) Subscription {
	return Subscription{
		Exchange: exchange,
		Pair:     pair,
		Channel:  ChannelOrderBook,
	}
}

// TradesSubs returns a subscription of channel 'trade'. Subscription's type will be set by Subscribe/Unsubscribe methods of
// WSConnection instance.
func TradesSubs(exchange, pair string) Subscription {
	return Subscription{
		Exchange: exchange,
		Pair:     pair,
		Channel:  ChannelTrades,
	}
}

// OrdersSubs returns a subscription of channel 'orders'. Subscription's type will be set by Subscribe/Unsubscribe methods of
// WSConnection instance.
func OrdersSubs() Subscription {
	return Subscription{
		Channel: ChannelOrders,
	}
}