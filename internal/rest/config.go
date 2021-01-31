package rest

type Config interface {
	PublicApiKey() string
	PrivateApiKey() string
}
