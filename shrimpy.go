package shrimpygo

import (
	"context"
	"fmt"
)

const (
	baseUrl   = "https://dev-api.shrimpy.io"
	wsBaseUrl = "wg://ws-feed.shrimpy.io"
	tokenPath = "/v1/ws/token"

	apiKeyHeader   = "DEV-SHRIMPY-API-KEY"
	apiNonceHeader = "DEV-SHRIMPY-API-NONCE"
	apiSigHeader   = "DEV-SHRIMPY-API-SIGNATURE"
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
	stream, err := createStream(ctx, shrimpy.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a websocket stream: %w", err)
	}

	return &WSConnection{stream: stream, ctx: ctx}, nil
}