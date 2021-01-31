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
		return fmt.Errorf("parser failed to decode data: %w", err), isPing
	}

	if data.Type != "" {
		if data.Type == "ping" {
			isPing = true
			return pingPong{Type: "pong", Data: data.Data}, isPing
		} else if data.Type == "error" {
			return fmt.Errorf("server error: code: %d, type: %s, message: %s",
				data.Code, data.Type, data.Message), isPing
		} else {
			return fmt.Errorf("parser: unknown data type: %v", data), isPing
		}
	}

	switch data.Channel {
	case "bbo", "orderbook":
		var priceData PriceData
		err = json.Unmarshal(b, &priceData)
		if err != nil {
			return fmt.Errorf("parser failed to decode data: expected to have price data from channel: %s, err: %w",
				data.Channel, err), isPing
		}
		return priceData, isPing
	}

	return fmt.Errorf("parser: unknown data type: channel: %s, data: %v", data.Channel, string(b)), isPing
}
