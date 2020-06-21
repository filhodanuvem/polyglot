package main

import (
	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/language"
	"github.com/filhodanuvem/polyglot/repository"
)

func main() {

	gh := github.Downloader{}
	path, err := gh.Download("https://github.com/filhodanuvem/gitql", "/Users/cloudson/sources/github/polyglot/temp")
	if err != nil {
		panic(err)
	}

	files := repository.GetFiles(path)
	var lang string
	for i := range files {
		lang, err = language.DetectByFile(files[i])
		if err != nil {
			panic(err)
		}

		println(files[i], lang)
	}
}
