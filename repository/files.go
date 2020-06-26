package repository

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func GetFiles(path string, l *log.Logger) []string {
	l.Printf("Listing files of %s\n", path)
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
