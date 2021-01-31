package ws

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

// Stream represents a websocket connection/stream which can be used to read from and write to the ws connection.
type Stream struct {
	conn *websocket.Conn
	// the stream receives user messages and ping replies through upstreamChan. note that we never close this channel
	// because the goroutine(s) writing to this channel is not its owner. eventually it get cleaned by the garbage collector.
	upstreamChan chan interface{}
	// server messages will be pushed to the data channel
	data chan []byte
	// error messages will be sent through the errors channel
	errors chan error
	// closing the errors chan needs synchronization since both downstream and upstream goroutines will write
	// to the same error channel.
	wg  sync.WaitGroup
	ctx context.Context
}

// CreateStream connects to the ws server and returns a Stream object containing the outbound channel.
// downstream (data) channel's capacity is set by throughput.
func CreateStream(ctx context.Context, throughput int, config StreamConfig) (*Stream, error) {
	wsConn, err := setupWSConnection(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup a ws connection: %w", err)
	}

	stream := &Stream{
		ctx:          ctx,
		conn:         wsConn,
		wg:           sync.WaitGroup{},
		errors:       make(chan error, 1),
		data:         make(chan []byte, throughput),
		upstreamChan: make(chan interface{}, 1),
	}

	stream.start()

	return stream, nil
}

func (s *Stream) start() {
	s.upstream()
	s.downstream()
	s.dispose()
}

// Send a message to websocket server.
func (s *Stream) Send(msg interface{}) {
	go func() {
		select {
		case <-s.ctx.Done():
		case s.upstreamChan <- msg:
		}
	}()
}

// DataChan returns a read-only channel used for reading from the ws connection. This channel will be closed in case of
// reading any error from the websocket connection.
func (s *Stream) DataChan() <-chan []byte {
	return s.data
}

// ErrorsChan returns a read-only channel used for reading ws connection errors. If a websocket read/write error occurs,
// it will be pushed to this channel and then the channel gets closed.
func (s *Stream) ErrorsChan() <-chan error {
	return s.errors
}

// downstream listens for incoming messages and push them to the data channel.
func (s *Stream) downstream() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(s.data)

		// reading the socket
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
			}
			_, msg, err := s.conn.ReadMessage()
			if err != nil {
				s.sendError(fmt.Errorf("websocket read error: %w", err))
				return
			}
			select {
			case <-s.ctx.Done():
				return
			case s.data <- msg:
			}
		}
	}()
}

// upstream listens for user (un)subscription messages, ping replies, etc; and sends them to the server.
func (s *Stream) upstream() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.ctx.Done():
				return
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
func (s *Stream) dispose() {
	go func() {
		s.wg.Wait()
		close(s.errors)
		_ = s.conn.Close()
	}()
}

func (s *Stream) sendError(err error) {
	select {
	case <-s.ctx.Done():
	case s.errors <- err:
	}
}
