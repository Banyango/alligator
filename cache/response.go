package cache

import (
	"bytes"
	"encoding/gob"
	"net/http"
)

type ResponseCache struct {
	Body      []byte
	HeaderMap http.Header
	Code      int
}

func ResponseCacheToBytes(cache *ResponseCache) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(cache)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ResponseCacheFromBytes(b []byte) (*ResponseCache, error) {
	var r ResponseCache
	enc := gob.NewDecoder(bytes.NewReader(b))
	err := enc.Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}