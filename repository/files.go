package repository

import (
	"os"
	"path/filepath"
)

func GetFiles(path string) []string {
	var files []string

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)

		return nil
	})

	return files
}
