package shrimpygo

import (
	"context"
	"fmt"
	"github.com/hedisam/shrimpygo/internal/rest"
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
// appropriate value as throughput to not miss any data. This is more important when you've subscribed to multiple
// channels of sequential data like order-book, especially in the times that market players are hyperactive.
func (cli *Client) Websocket(ctx context.Context, throughput int) (*WSConnection, error) {
	// connect to the ws server and create a ws stream
	stream, err := ws.CreateStream(ctx, throughput, cli.config)
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to create websocket stream: %w", err)
	}

	return &WSConnection{stream: stream, ctx: ctx, throughput: throughput}, nil
}

// SupportedExchanges (@PublicAPI) retrieves all Shrimpy supported exchanges and some basic information about each.
// no API keys are attached to the http request if you set freeApiCall which makes it a cost free request (rate limiting
// may be applied, check the docs @ https://developers.shrimpy.io/docs/#public)
func (cli *Client) SupportedExchanges(ctx context.Context, freeApiCall bool) ([]ExchangeInfo, error) {
	var exchanges []ExchangeInfo
	var cfg rest.Config = cli.config

	if freeApiCall {
		cfg = nil
	}

	err := rest.HttpGet(ctx, supportedExchanges, cfg, rest.NewDecoderFunc(&exchanges))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve the supported exchanges list: %w", err)
	}

	return exchanges, nil
}

// ExchangeAssets retrieves exchange asset information for a particular exchange.
// no API keys are attached to the http request if you set freeApiCall which makes it a cost free request (rate limiting
// may be applied, check the docs @ https://developers.shrimpy.io/docs/#public)
// Note that Shrimpy hosts a logo for each asset according to the table below:
// 32x32 	-> https://assets.shrimpy.io/cryptoicons/png/<id>.png
// 128x128 	-> https://assets.shrimpy.io/cryptoicons/png128/<id>.png
func (cli *Client) ExchangeAssets(ctx context.Context, exchange string, freeApiCall bool) ([]ExchangeAsset, error) {
	var assets []ExchangeAsset
	var urlPath = fmt.Sprintf(exchangeAssets, exchange)
	var cfg rest.Config = cli.config

	if freeApiCall {
		cfg = nil
	}

	err := rest.HttpGet(ctx, urlPath, cfg, rest.NewDecoderFunc(&assets))
	if err != nil {
		return nil, fmt.Errorf("shrimpygo failed to retrieve exchange assets list: %w", err)
	}

	return assets, nil
}