package rest

type Config interface {
	APIKey() string
	SecretKey() string
}
