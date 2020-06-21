package main

import (
	"fmt"
	"log"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
)

var limit = 4

func main() {
	repos, err := github.GetRepositories("filhodanuvem")
	if err != nil {
		panic(err)
	}

	stats := getStatisticsSync(repos)
	log.Println(stats.FirstLanguage())
}

func getStatisticsSync(repos []string) repository.Statistics {
	gh := github.Downloader{}
	var resultStats repository.Statistics
	c := 0
	for i := range repos {
		path, err := gh.Download(repos[i], "/Users/cloudson/sources/github/polyglot/temp")
		if err != nil {
			panic(err)
		}

		files := repository.GetFiles(path)
		stats, err := repository.GetStatistics(files)
		if err != nil {
			panic(err)
		}

		fmt.Println(stats)
		resultStats.Merge(&stats)
		c++
		if c == limit {
			break
		}
	}

	fmt.Println(resultStats)
	return resultStats
}
