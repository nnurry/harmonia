package utils

import "bytes"

type encoder interface {
	Encode(v any) error
}

func SerializeFromEncoder(encoder encoder, buf *bytes.Buffer, v any) ([]byte, error) {
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
