package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused session",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pomodoro resume")
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
