package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current pomodoro status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pomodoro status")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
