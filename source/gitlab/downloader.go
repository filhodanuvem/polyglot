package gitlab

import (
	"fmt"
	"strings"

	"github.com/filhodanuvem/polyglot/repository"
	log "github.com/sirupsen/logrus"
)

type Downloader struct {
}

func (d Downloader) Download(url, dest, defaultBranch string, l *log.Logger) (string, error) {
	l.WithFields(log.Fields{
		"repo": url,
		"dest": dest,
	}).Printf("Downloading repo")
	parts := strings.Split(url, "/")

	name := fmt.Sprintf("%s_%s", parts[len(parts)-2], parts[len(parts)-1])
	zipName := fmt.Sprintf("%s.zip", name)
	zipURL := fmt.Sprintf("%s/-/archive/%s/%s-%s.zip", url, defaultBranch, parts[4], defaultBranch)

	downloadedPath, err := repository.PrepareZIP(dest, name, zipURL, zipName, l)
	if err != nil {
		return "", err
	}

	return downloadedPath, err
}