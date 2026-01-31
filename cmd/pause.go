package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the current session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pomodoro pause")
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
