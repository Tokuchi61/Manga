package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const DefaultMaxBodyBytes int64 = 1 << 20 // 1 MiB

var ErrRequestBodyTooLarge = errors.New("request_body_too_large")

func DecodeJSON(r *http.Request, out any) error {
	payload, err := ReadBody(r, DefaultMaxBodyBytes)
	if err != nil {
		return err
	}
	return DecodeJSONBytes(payload, out)
}

func ReadBody(r *http.Request, maxBytes int64) ([]byte, error) {
	if r == nil || r.Body == nil {
		return nil, io.EOF
	}
	defer r.Body.Close()

	if maxBytes <= 0 {
		maxBytes = DefaultMaxBodyBytes
	}

	limited := io.LimitReader(r.Body, maxBytes+1)
	payload, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if int64(len(payload)) > maxBytes {
		return nil, ErrRequestBodyTooLarge
	}
	if len(bytes.TrimSpace(payload)) == 0 {
		return nil, io.EOF
	}

	return payload, nil
}

func DecodeJSONBytes(payload []byte, out any) error {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(out); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return fmt.Errorf("unexpected_trailing_json")
	}

	return nil
}