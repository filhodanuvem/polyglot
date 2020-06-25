package main

import (
	"log"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
)

var limit = 100
var tempPath = "/Users/cloudson/sources/github/polyglot/temp"

func main() {
	// log.SetOutput(ioutil.Discard)
	repos, err := github.GetRepositories("filhodanuvem")
	if err != nil {
		panic(err)
	}

	stats := getStatisticsAsync(repos)
	log.Printf("%+v", stats)
	log.Printf("%+v", stats.FirstLanguages(5))
}

func getStatisticsSync(repos []string) repository.Statistics {
	gh := github.Downloader{}
	var resultStats repository.Statistics
	c := 0
	for i := range repos {
		path, err := gh.Download(repos[i], tempPath)
		if err != nil {
			panic(err)
		}

		files := repository.GetFiles(path)
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

func getStatisticsAsync(repos []string) repository.Statistics {
	gh := github.Downloader{}
	statsChan := make(chan repository.Statistics, limit)
	done := make(chan bool, limit)
	count := 0

	for i := range repos {
		go func(repo string) {
			defer func() {
				done <- true
			}()
			path, err := gh.Download(repo, tempPath)
			if err != nil {
				log.Println(err)
				return
			}

			files := repository.GetFiles(path)
			stats, err := repository.GetStatistics(files)
			if err != nil {
				panic(err)
			}
			statsChan <- stats
		}(repos[i])
		count++
		if count == limit {
			break
		}
	}

	log.Println("Waiting for done")
	for range done {
		<-done
	}
	close(statsChan)

	log.Println("Waiting for statsChan")
	var resultStats repository.Statistics
	for range statsChan {
		stats := <-statsChan
		resultStats.Merge(&stats)
	}

	return resultStats
}
