package repository

import (
	"github.com/filhodanuvem/polyglot/source"
	log "github.com/sirupsen/logrus"
)

type Downloader interface {
	Download(repo source.ProviderRepo, dest string, l *log.Logger) (string, error)
}
