package keygen

import (
	"crypto/rand"
	"encoding/hex"
	"io"
)

func Keygen() string {
	key, err := GenerateKey()

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(key)
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, key)

	if err != nil {
		return nil, err
	}

	return key, nil
}

