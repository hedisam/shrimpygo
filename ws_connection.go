package shrimpygo

// WSConnection represents a shrimpy websocket connection
type WSConnection struct {
	stream *wsStream
}

// Send a message to the websocket server.
func (ws *WSConnection) Send(message interface{}) {
	ws.stream.send(message)
}
