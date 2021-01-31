package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/ws"
)

const (
	ChannelBBO       = "bbo"
	ChannelOrderBook = "orderbook"
	ChannelTrades    = "trade"
	ChannelOrders    = "orders"
)

type Shrimpy struct {
	config *Config
}

func NewShrimpyClient(cfx Config) (*Shrimpy, error) {
	if cfx.PublicApiKey() == "" || cfx.PrivateApiKey() == "" {
		return nil, fmt.Errorf("shrimpy config: invalid api/secret key")
	}

	return &Shrimpy{config: &cfx}, nil
}

// Websocket creates a websocket connection and returns a shrimpy WSConnection to interact with.
func (shrimpy *Shrimpy) Websocket(ctx context.Context) (*WSConnection, error) {
	// connect to the ws server and create a ws stream
	stream, err := ws.CreateStream(ctx, shrimpy.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a websocket stream: %w", err)
	}

	return &WSConnection{stream: stream, ctx: ctx}, nil
}
