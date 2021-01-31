package rest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Decoder func(reader io.Reader) error

func httpGet(ctx context.Context, path string, cfg Config, decoder Decoder) ([]byte, error) {
	signature, nonce, err := createSignature(cfg.PrivateApiKey(), path, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create api signature: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)
	if err != nil {
		return nil, fmt.Errorf("coudln't create http request: %w", err)
	}

	// signing the request
	req.Header.Add(apiKeyHeader, cfg.PublicApiKey())
	req.Header.Add(apiNonceHeader, fmt.Sprint(nonce))
	req.Header.Add(apiSigHeader, signature)

	req = req.WithContext(ctx)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unseccessful http request: status code: %d", resp.StatusCode)
	}

	err = decoder(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the http response body: %w", err)
	}

	return []byte{}, nil
}
