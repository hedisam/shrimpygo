package shrimpygo

import (
	"context"
)

// WSConnection represents a shrimpy websocket connection
type WSConnection struct {
	stream *wsStream
	ctx context.Context
}

// Send a message to the websocket server.
func (ws *WSConnection) Send(message interface{}) {
	ws.stream.send(message)
}

func (ws *WSConnection) Stream() <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		for {
			select {
			case <-ws.ctx.Done(): return
			case err := <-ws.stream.errorsChan():
				select {
				case <-ws.ctx.Done(): return
				case out <- err:
				}
			case rawMsg := <-ws.stream.dataChan():
				data, isPing := decode(rawMsg)
				if isPing {
					ws.Send(data)
					continue
				}
				select {
				case <-ws.ctx.Done(): return
				case out <- data:
				}
			}
		}
	}()
	return out
}
