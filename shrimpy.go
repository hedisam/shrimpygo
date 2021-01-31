package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/ws"
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
// throughput is just the downstream channel's capacity. the websocket connection replaces unread messages whenever
// new data comes in, so it's important to set an appropriate value to the throughput to not miss any data and this is
// more important when you're subscribed to (multiple) series of sequential data like orderbook.
func (shrimpy *Shrimpy) Websocket(ctx context.Context, throughput int) (*WSConnection, error) {
	// connect to the ws server and create a ws stream
	stream, err := ws.CreateStream(ctx, throughput, shrimpy.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a websocket stream: %w", err)
	}

	return &WSConnection{stream: stream, ctx: ctx, throughput: throughput}, nil
}
