package cmd

import (
	"fmt"
	"os"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/repository"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var limitRepos = 100
var limitChannels = 30
var tempPath = "/tmp/polyglot"

var logLevels = map[string]log.Level{
	"debug":   log.DebugLevel,
	"info":    log.InfoLevel,
	"warning": log.WarnLevel,
	"error":   log.ErrorLevel,
	"fatal":   log.FatalLevel,
}

func Run(cmd *cobra.Command, args []string) {
	l := log.New()
	logVerbosity, _ := cmd.Flags().GetString("log")
	if level, ok := logLevels[logVerbosity]; ok {
		l.SetLevel(level)
	}

	l.SetOutput(os.Stdout)
	outputFile, _ := cmd.Flags().GetString("output")
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		l.SetOutput(file)
	}

	username, err := cmd.Flags().GetString("username")
	if err != nil {
		panic(err)
	}
	repos, err := github.GetRepositories(username)
	if err != nil {
		l.Println(err)
	}
	stats := getStatisticsSync(repos, l)
	fmt.Printf("First 5 languages\n%+v", stats.FirstLanguages(5))
}

func getStatisticsAsync(repos []string, l *log.Logger) repository.Statistics {
	statsChan := make(chan repository.Statistics, limitRepos)
	terminated := make(chan bool, limitChannels)
	count := 0

	for i := range repos {
		go func(repo string) {
			defer func() {
				terminated <- true
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

	l.Println(">>>>>> Waiting for terminated")
	for i := 0; i < count; i++ {
		select {
		case <-terminated:
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
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		os.MkdirAll(tempPath, os.ModePerm)
	}
	gh := github.Downloader{}
	path, err := gh.Download(repo, tempPath, l)
	if err != nil {
		l.Error(err)
	}

	files := repository.GetFiles(path, l)
	stats, err := repository.GetStatistics(files)

	return stats, err
}
