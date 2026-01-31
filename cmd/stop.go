package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pomodoro stop")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
