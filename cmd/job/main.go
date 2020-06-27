package main

import (
	"fmt"
	"os"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
	log "github.com/sirupsen/logrus"
)

var limitRepos = 100
var limitChannels = 30
var tempPath = "/Users/cloudson/sources/github/polyglot/temp"

func main() {
	l := log.New()
	// l.SetLevel(log.WarnLevel)
	l.SetOutput(os.Stdout)
	repos, err := github.GetRepositories("filhodanuvem")
	if err != nil {
		l.Println(err)
	}

	stats := getStatisticsAsync(repos, l)
	fmt.Printf("First 5 languages\n%+v", stats.FirstLanguages(5))
}

func getStatisticsAsync(repos []string, l *log.Logger) repository.Statistics {
	statsChan := make(chan repository.Statistics, limitRepos)
	done := make(chan bool, limitChannels)
	count := 0

	for i := range repos {
		go func(repo string) {
			defer func() {
				done <- true
			}()
			stats, err := getStatsFromRepo(repo, tempPath, l)
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

	l.Println(">>>>>> Waiting for done")
	for i := 0; i < count; i++ {
		select {
		case <-done:
		}
	}
	close(statsChan)

	l.Println(">>>>>>> Waiting for statsChan")
	var resultStats repository.Statistics
	for range statsChan {
		stats := <-statsChan
		resultStats.Merge(&stats)
		l.Infof(">>Merge %+v making %+v", stats, resultStats)
	}

	return resultStats
}

func getStatisticsSync(repos []string, l *log.Logger) repository.Statistics {
	var resultStats repository.Statistics
	c := 0
	for i := range repos {
		stats, err := getStatsFromRepo(repos[i], tempPath, l)
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

func getStatsFromRepo(repo, tempPath string, l *log.Logger) (repository.Statistics, error) {
	gh := github.Downloader{}
	path, err := gh.Download(repo, tempPath, l)
	if err != nil {
		l.Error(err)
	}

	files := repository.GetFiles(path, l)
	stats, err := repository.GetStatistics(files)

	return stats, err
}
