package cmd

import (
	"github.com/bata94/reqlab/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Show the largest files in the given path.",
	Long:  `Quickly scan a directory and find large files.`,
	Run: func(cmd *cobra.Command, args []string) {
		tui.MainView()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
