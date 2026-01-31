package cmd

import (
	"fmt"
	"time"

	"pomodoro/internal/state"

	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused session",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("load state: %w", err)
		}

		if s == nil {
			fmt.Println("No session found.")
			return nil
		}

		if !s.Running {
			fmt.Println("🛑 No running session to resume.")
			fmt.Println("Hint: pomodoro start")
			return nil
		}

		if !s.Paused {
			fmt.Println("Session is not paused.")
			return nil
		}

		remaining := s.RemainingSeconds
		if remaining < 0 {
			remaining = 0
		}

		now := time.Now().UTC()
		s.Paused = false
		s.PausedAt = nil
		s.RemainingSeconds = remaining
		s.StartedAt = now.Format(time.RFC3339)
		s.EndsAt = now.Add(time.Duration(remaining) * time.Second).Format(time.RFC3339)

		if err := state.Save(s); err != nil {
			return fmt.Errorf("save state: %w", err)
		}

		fmt.Println("▶️ Pomodoro resumed.")
		fmt.Printf("Remaining: %s\n", formatMMSS(remaining))
		fmt.Printf("Ends at: %s\n", now.Add(time.Duration(remaining)*time.Second).Local().Format("15:04:05"))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}
