package shrimpygo

// Config holds api keys. You can find your API keys from shrimpy developers dashboard.
type Config struct {
	// Public Key
	PublicKey string
	// Private Key
	PrivateKey string
}

func (cfg *Config) PublicApiKey() string {
	return cfg.PublicKey
}

func (cfg *Config) PrivateApiKey() string {
	return cfg.PrivateKey
}
