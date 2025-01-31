package server

import (
	"bytes"
	"encoding/json"
)

func encodeJSON[T any](val T) ([]byte, error) {
	var b bytes.Buffer

	err := json.NewEncoder(&b).Encode(val)
	return b.Bytes(), err
}

func decodeJSON[T any](b *bytes.Buffer) (T, error) {
	var out T
	err := json.NewDecoder(b).Decode(&out)
	return out, err
}
