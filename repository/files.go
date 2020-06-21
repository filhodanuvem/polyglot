package repository

import (
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

func GetFiles(path string) []string {
	log.Infof("Listing files of %s", path)
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
