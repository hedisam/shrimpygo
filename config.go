package shrimpygo

type shrimpyConfig struct {
	apiKey string
	secretKey string
}

func (cfg *shrimpyConfig) APIKey() string {
	return cfg.apiKey
}

func (cfg *shrimpyConfig) SecretKey() string {
	return cfg.secretKey
}