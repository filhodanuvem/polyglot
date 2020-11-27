package github

import (
	"fmt"
	"strings"

	"github.com/filhodanuvem/polyglot/repository"
	"github.com/filhodanuvem/polyglot/source"
	log "github.com/sirupsen/logrus"
)

type Downloader struct {
}

func (d Downloader) Download(repo source.ProviderRepo, dest string, l *log.Logger) (string, error) {
	l.WithFields(log.Fields{
		"repo": repo,
		"dest": dest,
	}).Printf("Downloading repo")
	parts := strings.Split(repo.URL, "/")
	name := fmt.Sprintf("%s_%s", parts[len(parts)-2], parts[len(parts)-1])
	zipName := fmt.Sprintf("%s.zip", name)
	zipURL := fmt.Sprintf("%s/archive/%s.zip", repo.URL, repo.DefaultBranch)

	downloadedPath, err := repository.PrepareZIP(dest, name, zipURL, zipName, l)
	if err != nil {
		return "", err
	}

	return downloadedPath, err
}
