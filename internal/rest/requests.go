package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)


func HttpGet(ctx context.Context, path string, cfg Config, decoder Decoder) error {
	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)
	if err != nil {
		return fmt.Errorf("coudln't create http request: %w", err)
	}

	// nil cfg is only acceptable for public api calls
	if cfg != nil {
		signature, nonce, err := createSignature(cfg.PrivateApiKey(), path, http.MethodGet, nil)
		if err != nil {
			return fmt.Errorf("failed to create api signature: %w", err)
		}

		// api keys & signing the request
		req.Header.Add(apiKeyHeader, cfg.PublicApiKey())
		req.Header.Add(apiNonceHeader, fmt.Sprint(nonce))
		req.Header.Add(apiSigHeader, signature)
	}

	req = req.WithContext(ctx)

	// todo: read timeout from the config object
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unseccessful http request: status code: %d", resp.StatusCode)
	}

	err = decoder(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read the http response body: %w", err)
	}

	return nil
}
