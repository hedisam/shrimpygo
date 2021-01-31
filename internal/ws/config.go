package ws

type StreamConfig interface {
	APIKey() string
	SecretKey() string
}
