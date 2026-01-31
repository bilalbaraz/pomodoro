package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"pomodoro/internal/state"

	"github.com/spf13/cobra"
)

var statusJSON bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current pomodoro status",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := state.Load()
		if err != nil {
			return fmt.Errorf("load state: %w", err)
		}

		if s == nil {
			return printStatusJSONOrText(
				map[string]any{
					"status":  "stopped",
					"message": "No active session. Run: pomodoro start",
				},
				func() {
					fmt.Println("No active session. Run: pomodoro start")
				},
			)
		}

		if !s.Running {
			return printStatusJSONOrText(
				map[string]any{
					"status":  "stopped",
					"message": "No active session.",
				},
				func() {
					fmt.Println("🛑 Status: STOPPED")
					fmt.Println("No active session.")
					fmt.Println("Hint: pomodoro start")
				},
			)
		}

		endsAt, err := time.Parse(time.RFC3339, s.EndsAt)
		if err != nil {
			return fmt.Errorf("parse ends_at: %w", err)
		}

		now := time.Now().UTC()
		if endsAt.Before(now) {
			s.Running = false
			if saveErr := state.Save(s); saveErr != nil {
				return fmt.Errorf("save finished state: %w", saveErr)
			}
			return printStatusJSONOrText(
				map[string]any{
					"status":  "stopped",
					"message": "Session already finished.",
				},
				func() {
					fmt.Println("Session already finished.")
				},
			)
		}

		statusJSONLabel := "running"
		if s.Paused {
			statusJSONLabel = "paused"
		}

		remaining := s.RemainingSeconds
		if !s.Paused {
			remaining = int(endsAt.Sub(now).Seconds())
			if remaining < 0 {
				remaining = 0
			}
		}

		sessionStr := fmt.Sprintf("%d/%d", s.SessionIndex, s.SessionsTotal)

		jsonPayload := map[string]any{
			"status":            statusJSONLabel,
			"mode":              s.Mode,
			"remaining_seconds": remaining,
			"task":              s.Task,
			"session":           sessionStr,
		}

		return printStatusJSONOrText(
			jsonPayload,
			func() {
				if s.Paused {
					fmt.Println("⏸️ Status: PAUSED")
				} else {
					fmt.Println("🍅 Status: RUNNING")
				}
				fmt.Printf("Mode: %s (session %s)\n", s.Mode, sessionStr)
				if s.Task != "" {
					fmt.Printf("Task: %s\n", s.Task)
				}
				fmt.Printf("Remaining: %s\n", formatMMSS(remaining))
				fmt.Printf("Ends at: %s\n", endsAt.Local().Format("15:04:05"))
				if s.Paused {
					fmt.Println("Hint: pomodoro resume")
				}
			},
		)
	},
}

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "Output status as JSON")
	rootCmd.AddCommand(statusCmd)
}

func printStatusJSONOrText(payload map[string]any, text func()) error {
	if statusJSON {
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return fmt.Errorf("encode status json: %w", err)
		}
		fmt.Fprintln(os.Stdout, string(data))
		return nil
	}

	text()
	return nil
}
