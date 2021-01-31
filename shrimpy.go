package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/ws"
)

const (
	ChanBBO       = "bbo"
	ChanOrderBook = "orderbook"
	ChanTrades    = "trades"
	ChanOrders    = "orders"
)

type Shrimpy struct {
	config *shrimpyConfig
}

func NewShrimpyClient(apiKey, secretKey string) (*Shrimpy, error) {
	if apiKey == "" || secretKey == "" {
		return nil, fmt.Errorf("shrimpy: invalid api/secret key")
	}

	shConfig := &shrimpyConfig{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
	return &Shrimpy{config: shConfig}, nil
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
