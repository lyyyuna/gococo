package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gococo",
	Short: "gococo is another Golang coverage collection tool",
	Long:  "gococo is a free tool for Golang Coverage Collection",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
	}
}
