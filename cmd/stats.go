package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show pomodoro statistics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pomodoro stats")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
