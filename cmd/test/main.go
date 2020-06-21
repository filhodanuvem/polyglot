package main

import (
	"log"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
)

func main() {

	gh := github.Downloader{}
	path, err := gh.Download("https://github.com/filhodanuvem/gitql", "/Users/cloudson/sources/github/polyglot/temp")
	if err != nil {
		panic(err)
	}

	files := repository.GetFiles(path)
	stats, err := repository.GetStatistics(files)
	if err != nil {
		panic(err)
	}

	log.Println(stats.FirstLanguage())
}
