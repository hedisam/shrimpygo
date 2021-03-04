package rest

import (
	"encoding/json"
	"fmt"
	"io"
)

// Decoder is used to decode http response body.
type Decoder func(reader io.Reader) error

func NewDecoderFunc(v interface{}) Decoder {
	return func(reader io.Reader) error {
		err := json.NewDecoder(reader).Decode(v)
		if err != nil {
			return fmt.Errorf("failed to decode http response: %w", err)
		}
		return nil
	}
}
