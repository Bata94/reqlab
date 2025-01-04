package cmd

import (
	"github.com/bata94/reqlab/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the TUI",
	Long:  "Launch the TUI, this is the main use of the App",
	Run: func(cmd *cobra.Command, args []string) {
		tui.MainView()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
