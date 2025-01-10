package cmd

import (
	"github.com/bata94/reqlab/internal/loadtest"
	"github.com/spf13/cobra"
)

var ltCmd = &cobra.Command{
	Use:     "loadtest",
	Aliases: []string{"lt"},
	Short:   "Start a loadtest",
	Long:    "Start a loadtest",
	Run: func(cmd *cobra.Command, args []string) {
		loadtest.Main()
	},
}

func init() {
	rootCmd.AddCommand(ltCmd)
}
