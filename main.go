package main

import (
	"fmt"
	"os"

	"github.com/filhodanuvem/polyglot/cmd"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Polyglot",
	Short: "Polyglot tells you the (programming) languages that you speak",
	Run:   cmd.Run,
}

func main() {
	rootCmd.PersistentFlags().StringP("log", "l", "fatal", "Log verbosity, options [debug, info, warning, error, fatal]")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Path to log in a file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
