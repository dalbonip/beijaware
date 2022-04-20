package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

//959148601e08d6cc80816964006a15ed54911fd32c3a77f7baaf5d74d4f895c0
func Keygen() string {
	key, err := GenerateKey()

	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(key)
}

func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	fmt.Println(key)
	_, err := io.ReadFull(rand.Reader, key)

	if err != nil {
		return nil, err
	}

	key = []byte("headmindpartnersredteam1337@1337")

	return key, nil
}
func main() {
	key := Keygen()
	fmt.Println(key)
}
