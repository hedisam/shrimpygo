package shrimpygo

import (
	"context"
	"github.com/hedisam/shrimpygo/internal/ws"
)

// WSConnection represents a shrimpy websocket connection
type WSConnection struct {
	stream *ws.Stream
	ctx    context.Context
}

// Subscribe to a channel on this ws connection.
func (ws *WSConnection) Subscribe(channel, exchange, pair string) {
	ws.stream.Send(Subscription{
		Type:     "subscribe",
		Exchange: exchange,
		Pair:     pair,
		Channel:  channel,
	})
}

// Unsubscribe from a channel on this ws connection.
func (ws *WSConnection) Unsubscribe(channel, exchange, pair string) {
	ws.stream.Send(Subscription{
		Type:     "unsubscribe",
		Exchange: exchange,
		Pair:     pair,
		Channel:  channel,
	})
}

func (ws *WSConnection) Stream() <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		for {
			select {
			case <-ws.ctx.Done():
				return
			case err := <-ws.stream.ErrorsChan():
				select {
				case <-ws.ctx.Done():
					return
				case out <- err:
				}
			case rawMsg := <-ws.stream.DataChan():
				data, isPing := decode(rawMsg)
				if isPing {
					ws.stream.Send(data)
					continue
				}
				select {
				case <-ws.ctx.Done():
					return
				case out <- data:
				}
			}
		}
	}()
	return out
}
