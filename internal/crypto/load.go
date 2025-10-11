package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadRSAKey(filename string) *rsa.PrivateKey {
	data, err := os.ReadFile(string(filename))
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		panic(fmt.Errorf("no PEM formatted block found"))
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	return key
}
