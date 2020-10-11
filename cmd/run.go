package cmd

import (
	"fmt"
	"os"

	"github.com/filhodanuvem/polyglot/github"
	"github.com/filhodanuvem/polyglot/server"
	"github.com/filhodanuvem/polyglot/stats"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var limitRepos = 100
var limitChannels = 30

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
	tempPath, _ := cmd.Flags().GetString("path")
	useServer, _ := cmd.Flags().GetBool("server")

	if useServer {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetString("port")
		server.Serve(host, port, tempPath)
	} else {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			panic(err)
		}
		if username == "" {
			fmt.Println("required flag(s) \"username\" not set")
			cmd.Help()
			os.Exit(1)
		}
		repos, err := github.GetRepositories(username)
		if err != nil {
			l.Println(err)
		}
		stats := stats.GetStatisticsAsync(tempPath, repos, l)
		fmt.Printf("First 5 languages\n%+v\n", stats.FirstLanguages(5))

	}
}
