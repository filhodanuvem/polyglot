package github

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"archive/zip"

	log "github.com/sirupsen/logrus"
)

type Downloader struct {
}

func (d *Downloader) Download(url, dest string, l *log.Logger) (string, error) {
	l.WithFields(log.Fields{
		"repo": url,
		"dest": dest,
	}).Printf("Downloading repo")
	parts := strings.Split(url, "/")
	name := fmt.Sprintf("%s_%s", parts[len(parts)-2], parts[len(parts)-1])
	zipName := fmt.Sprintf("%s.zip", name)
	zipURL := fmt.Sprintf("%s/archive/master.zip", url)

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

	return downloadedPath, err
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
