package shrimpygo

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type wsStream struct {
	conn *websocket.Conn
	// the stream receives user messages and ping replies through upstreamChan. note that we never close this channel
	// because the goroutine(s) writing to this channel is not its owner. eventually it get cleaned by the garbage collector.
	upstreamChan chan interface{}
	// server messages will be pushed to the data channel
	data   chan []byte
	// error messages will be sent through the errors channel
	errors chan error
	// closing the errors chan needs synchronization since both downstream and upstream goroutines will write
	// to the same error channel.
	wg sync.WaitGroup
	ctx context.Context
}

// createStream connects to the ws server and returns a wsStream object containing the outbound channel.
func createStream(ctx context.Context, config *shrimpyConfig) (*wsStream, error) {
	wsConn, err := setupWSConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup a ws connection: %w", err)
	}

	stream := &wsStream{
		conn:         wsConn,
		data:         make(chan []byte),
		errors:       make(chan error),
		upstreamChan: make(chan interface{}),
		wg:           sync.WaitGroup{},
		ctx:          ctx,
	}

	stream.start()

	return stream, nil
}

func (s *wsStream) start() {
	s.upstream()
	s.downstream()
	s.dispose()
}

func (s *wsStream) send(msg interface{}) {
	go func() {
		select {
		case <-s.ctx.Done():
		case s.upstreamChan <- msg:
		}
	}()
}

// downstream listens for incoming messages and push them to the data channel.
func (s *wsStream) downstream() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(s.data)

		// reading the socket
		for {
			select {
			case <-s.ctx.Done(): return
			default:
			}
			_, msg, err := s.conn.ReadMessage()
			if err != nil {
				s.sendError(fmt.Errorf("websocket read error: %w", err))
				return
			}
			select {
			case <-s.ctx.Done(): return
			case s.data <- msg:
			}
		}
	}()
}

// upstream listens for user (un)subscription messages, ping replies, etc; and sends them to the server.
func (s *wsStream) upstream() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ctx.Done(): return
			case msg, ok := <-s.upstreamChan:
				if !ok {
					return
				}
				err := s.conn.WriteJSON(msg)
				if err != nil {
					s.sendError(fmt.Errorf("failed to send message to the ws server: %w", err))
					return
				}
			}
		}
	}()
}

// dispose waits for upstream and downstream goroutines to get done then closes the shared errors channel and
// the websocket connection (which could be closed already).
func (s *wsStream) dispose() {
	go func() {
		s.wg.Wait()
		close(s.errors)
		_ = s.conn.Close()
	}()
}

func (s *wsStream) sendError(err error) {
	select {
	case <-s.ctx.Done():
	case s.errors <- err:
	}
}


