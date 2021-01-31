package shrimpygo

import (
	"context"
	"github.com/hedisam/shrimpygo/internal/ws"
)

// WSConnection represents a shrimpy websocket connection
type WSConnection struct {
	stream     *ws.Stream
	ctx        context.Context
	throughput int
}

// Subscribe to a channel on this ws connection.
func (ws *WSConnection) Subscribe(subscriptions ...Subscription) {
	for _, subs := range subscriptions {
		s := subs
		s.Type = "subscribe"
		ws.stream.Send(s)
	}
}

// Unsubscribe from a channel on this ws connection.
func (ws *WSConnection) Unsubscribe(unSubscriptions ...Subscription) {
	for _, unSub := range unSubscriptions {
		u := unSub
		u.Type = "unsubscribe"
		ws.stream.Send(u)
	}
}

// Stream returns a read-only channel which can be used to read websocket messages along with all kind of errors that
// could happen while reading from or creating the ws connection.
func (ws *WSConnection) Stream() <-chan interface{} {
	out := make(chan interface{}, ws.throughput)
	go func() {
		defer close(out)

		for {
			select {
			case <-ws.ctx.Done():
				return
			case err, ok := <-ws.stream.ErrorsChan():
				if !ok {return}
				select {
				case <-ws.ctx.Done():
					return
				case out <- err:
				}
			case rawMsg, ok := <-ws.stream.DataChan():
				if !ok {return}
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
