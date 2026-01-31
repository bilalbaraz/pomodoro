package cmd

import (
	"fmt"
	"time"

	"pomodoro/internal/state"

	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the current session",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("load state: %w", err)
		}

		if s == nil {
			fmt.Println("No active session to pause.")
			return nil
		}

		if !s.Running {
			fmt.Println("🛑 No running session.")
			fmt.Println("Hint: pomodoro start")
			return nil
		}

		if s.Paused {
			fmt.Println("Session already paused.")
			return nil
		}

		endsAt, err := time.Parse(time.RFC3339, s.EndsAt)
		if err != nil {
			return fmt.Errorf("parse ends_at: %w", err)
		}

		now := time.Now().UTC()
		remaining := int(endsAt.Sub(now).Seconds())
		if remaining < 0 {
			remaining = 0
		}

		pausedAt := now.Format(time.RFC3339)
		s.Paused = true
		s.PausedAt = &pausedAt
		s.RemainingSeconds = remaining

		if err := state.Save(s); err != nil {
			return fmt.Errorf("save state: %w", err)
		}

		fmt.Println("⏸️ Pomodoro paused.")
		fmt.Printf("Remaining: %s\n", formatMMSS(remaining))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
}
