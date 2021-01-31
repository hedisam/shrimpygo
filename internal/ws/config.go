package ws

type StreamConfig interface {
	PublicApiKey() string
	PrivateApiKey() string
}
