package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "oebb-cli",
	Short: "A command line client for the ÖBB Tickets API",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func Execute() {
	searchCmd.Flags().IntP("results", "n", 5, "Number of search results to display")
	searchCmd.Flags().StringP("time", "t", "", "Departure time")
	rootCmd.AddCommand(searchCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
