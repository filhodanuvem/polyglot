package repository

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func GetFiles(path string, l *log.Logger) []string {
	l.WithFields(log.Fields{
		"path": path,
	}).Printf("Listing files of repo")
	var files []string

	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		files = append(files, path)

		return nil
	})

	return files
}
