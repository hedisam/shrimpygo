package ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/hedisam/shrimpygo/internal/rest"
)

func setupWSConnection(ctx context.Context, config StreamConfig) (*websocket.Conn, error) {
	// the ws server requires a valid token
	token, err := rest.Token(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("couldn't get a websocket token: %w", err)
	}

	// connecting to the server
	url := fmt.Sprintf("%s?token=%s", wsBaseUrl, token)
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket: %w", err)
	}
	return conn, nil
}
