package env

import (
	"encoding/base64"
)

type Base64EncodedBytes []byte

func (b *Base64EncodedBytes) UnmarshalText(text []byte) error {
	out, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}

	*b = out
	return nil

}
