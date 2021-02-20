package stats

import (
	"os"

	"github.com/filhodanuvem/polyglot/repository"
	"github.com/filhodanuvem/polyglot/source"
	"github.com/filhodanuvem/polyglot/source/github"
	"github.com/filhodanuvem/polyglot/source/gitlab"
	log "github.com/sirupsen/logrus"
)

var limitRepos = 100
var limitChannels = 30

func GetStatisticsAsync(tempPath, provider string, repos []source.ProviderRepo, l *log.Logger) repository.Statistics {
	statsChan := make(chan repository.Statistics, limitRepos)
	count := 0

	for i := range repos {
		go func(repo source.ProviderRepo) {
			stats, err := getStatsFromRepo(repo, tempPath, provider, l)
			if err != nil {
				l.Error(err)
				return
			}
			statsChan <- stats
		}(repos[i])
		count++
		if count == limitRepos {
			break
		}
	}

	var resultStats repository.Statistics
	for range repos {
		select {
		case stats := <-statsChan:
			resultStats.Merge(&stats)
		}
	}

	return resultStats
}

func GetStatisticsSync(tempPath, provider string, repos []source.ProviderRepo, l *log.Logger) repository.Statistics {
	var resultStats repository.Statistics
	c := 0
	for i := range repos {
		stats, err := getStatsFromRepo(repos[i], tempPath, provider, l)
		if err != nil {
			l.Error(err)
			continue
		}
		resultStats.Merge(&stats)
		c++
		if c == limitRepos {
			break
		}
	}

	return resultStats
}

func getStatsFromRepo(repo source.ProviderRepo, tempPath, provider string, l *log.Logger) (repository.Statistics, error) {
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		os.MkdirAll(tempPath, os.ModePerm)
	}

	var path string
	var err error

	var downloader interface{} = github.Downloader{}
	if provider == "gitlab" {
		downloader = gitlab.Downloader{}
	}

	tDownloader := downloader.(repository.Downloader)

	path, err = tDownloader.Download(repo, tempPath, l)
	if err != nil {
		l.Error(err)
	}

	files := repository.GetFiles(path, l)
	stats, err := repository.GetStatistics(files)

	return stats, err
}
