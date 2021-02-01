package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/ws"
)

type Client struct {
	config *Config
}

func NewClient(cfx Config) (*Client, error) {
	if cfx.PublicApiKey() == "" || cfx.PrivateApiKey() == "" {
		return nil, fmt.Errorf("shrimpy config: invalid api/secret key")
	}

	return &Client{config: &cfx}, nil
}

// Websocket creates a websocket connection and returns a shrimpy WSConnection to interact with.
// throughput is just the downstream channel's capacity. Although in normal circumstances a small throughput works but
// since the websocket connection replaces unread messages whenever new data comes in, it's important to set an
// appropriate value as throughput to not miss any data. This is more important when you're subscribed to multiple
// channels of sequential data like order-book, especially in the times that market players are hyperactive.
func (shrimpy *Client) Websocket(ctx context.Context, throughput int) (*WSConnection, error) {
	// connect to the ws server and create a ws stream
	stream, err := ws.CreateStream(ctx, throughput, shrimpy.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a websocket stream: %w", err)
	}

	return &WSConnection{stream: stream, ctx: ctx, throughput: throughput}, nil
}
