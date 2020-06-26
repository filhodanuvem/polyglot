package main

import (
	"os"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
	log "github.com/sirupsen/logrus"
)

var limit = 100
var tempPath = "/Users/cloudson/sources/github/polyglot/temp"

func main() {
	l := log.New()
	l.SetLevel(log.WarnLevel)
	l.SetOutput(os.Stdout)
	repos, err := github.GetRepositories("filhodanuvem")
	if err != nil {
		l.Println(err)
	}

	stats := getStatisticsAsync(repos, l)
	l.Printf("First 5 languages\n%+v", stats.FirstLanguages(25))
}

func getStatisticsSync(repos []string, l *log.Logger) repository.Statistics {
	gh := github.Downloader{}
	var resultStats repository.Statistics
	c := 0
	for i := range repos {
		path, err := gh.Download(repos[i], tempPath, l)
		if err != nil {
			panic(err)
		}

		files := repository.GetFiles(path, l)
		stats, err := repository.GetStatistics(files)
		if err != nil {
			panic(err)
		}
		resultStats.Merge(&stats)
		c++
		if c == limit {
			break
		}
	}

	return resultStats
}

func getStatisticsAsync(repos []string, l *log.Logger) repository.Statistics {
	gh := github.Downloader{}
	statsChan := make(chan repository.Statistics, limit)
	done := make(chan bool, limit)
	count := 0

	for i := range repos {
		go func(repo string) {
			defer func() {
				done <- true
			}()
			path, err := gh.Download(repo, tempPath, l)
			if err != nil {
				l.Error(err)
				return
			}

			files := repository.GetFiles(path, l)
			stats, err := repository.GetStatistics(files)
			if err != nil {
				l.Error(err)
			}
			statsChan <- stats
		}(repos[i])
		count++
		if count == limit {
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
	}

	return resultStats
}
