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
	rootCmd.Flags().StringP("username", "u", "", "Username")
	rootCmd.Flags().StringP("path", "p", "/tmp/polyglot", "Path where to download the repositories")
	rootCmd.Flags().StringP("provider", "", "github", "Repository Provider, options [github, gitlab]")
	rootCmd.PersistentFlags().StringP("log", "l", "fatal", "Log verbosity, options [debug, info, warning, error, fatal]")
	rootCmd.PersistentFlags().StringP("output", "o", "", "Path to log in a file")
	rootCmd.Flags().BoolP("server", "s", false, "Run polyglot API Server")
	rootCmd.PersistentFlags().StringP("host", "", "127.0.0.1", "IP address for the server")
	rootCmd.PersistentFlags().StringP("port", "", "8080", "Port for the server")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
