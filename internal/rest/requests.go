package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func get(ctx context.Context, path string, cfg Config) ([]byte, error) {
	signature, nonce, err := createSignature(cfg.SecretKey(), path, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create api signature: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, baseUrl+path, nil)
	if err != nil {
		return nil, fmt.Errorf("coudln't create http request: %w", err)
	}

	// signing the request
	req.Header.Add(apiKeyHeader, cfg.APIKey())
	req.Header.Add(apiNonceHeader, fmt.Sprint(nonce))
	req.Header.Add(apiSigHeader, signature)

	// todo: use the context
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request not accepted: status code: %d", resp.StatusCode)
	}

	// Todo: don't use ioutil RealAll. Response body might be large.
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read the response body: %w", err)
	}

	return b, nil
}
