package shrimpygo

import (
	"encoding/json"
	"fmt"
)

func decode(b []byte) (interface{}, bool) {
	var isPing bool
	var data unknownData
	err := json.Unmarshal(b, &data)
	if err != nil {
		return fmt.Errorf("failed to decode json data: %w", err), isPing
	}

	if data.Type != "" {
		if data.Type == "ping" {
			isPing = true
			return &pingPong{Type: "pong", Data: data.Data}, isPing
		} else if data.Type == "error" {
			return fmt.Errorf("server error: code: %d, type: %s, message: %s",
				data.Code, data.Type, data.Message), isPing
		} else {
			return fmt.Errorf("decode: unknown data type: %v", data), isPing
		}
	}

	switch data.Channel {
	case ChannelBBO, ChannelOrderBook:
		var orderBook OrderBookInfo
		err = json.Unmarshal(b, &orderBook)
		if err != nil {
			return fmt.Errorf("decode failed to decode data: expected to have orderbook/bbo data from channel: %s, err: %w",
				data.Channel, err), isPing
		}
		return &orderBook, isPing
	case ChannelTrades:
		var trades Trades
		err = json.Unmarshal(b, &trades)
		if err != nil {
			return fmt.Errorf("decode failed to decode data: expected to have trades data from channel: %s, err: %w",
				data.Channel, err), isPing
		}
		return &trades, isPing
	case ChannelOrders:
		var orders Orders
		err = json.Unmarshal(b, &orders)
		if err != nil {
			return fmt.Errorf("decode failed to decode data: expected to have orders data from channel: %s, err: %w",
				data.Channel, err), isPing
		}
		return &orders, isPing
	}

	return fmt.Errorf("decode: unknown data type: channel: %s, data: %v", data.Channel, string(b)), isPing
}
