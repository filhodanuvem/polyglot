package repository

import (
	"archive/zip"
	"io"
	"net/http"
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

func PrepareZIP(dest, name, zipURL, zipName string, l *log.Logger) (string, error) {
	resp, err := http.Get(zipURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	path := filepath.Join(dest, zipName)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	downloadedPath := filepath.Join(dest, name)

	err = unzip(path, downloadedPath, l)
	if err != nil {
		return "", err
	}

	return downloadedPath, nil
}

func unzip(path, dest string, l *log.Logger) error {
	l.WithFields(log.Fields{
		"from": path,
		"dest": dest,
	}).Printf("Unzipping repo")
	reader, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer reader.Close()
	defer os.Remove(path)

	for i := range reader.File {
		file := reader.File[i]
		zippedFile, err := file.Open()
		if err != nil {
			return err
		}

		extractedFilePath := filepath.Join(
			dest,
			file.Name,
		)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
			continue
		}

		outputFile, err := os.OpenFile(
			extractedFilePath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			file.Mode(),
		)
		if err != nil {
			return err
		}
		_, err = io.Copy(outputFile, zippedFile)
		if err != nil {
			return err
		}
		outputFile.Close()
		zippedFile.Close()
	}

	return nil
}
