package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a pomodoro session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pomodoro start")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
