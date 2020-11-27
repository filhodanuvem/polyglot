package repository

import (
	log "github.com/sirupsen/logrus"
)

type Downloader interface {
	Download(url, dest, defaultBranch string, l *log.Logger) (string, error)
}
