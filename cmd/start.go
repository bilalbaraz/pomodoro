package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultWorkMinutes  = 25
	defaultBreakMinutes = 5
)

var (
	workMinutes  int
	breakMinutes int
	taskName     string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a pomodoro session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workMinutes <= 0 {
			return fmt.Errorf("work minutes must be greater than 0")
		}
		if breakMinutes <= 0 {
			return fmt.Errorf("break minutes must be greater than 0")
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		printStart(taskName)
		runCountdown(ctx, time.Duration(workMinutes)*time.Minute)

		fmt.Println("✅ Pomodoro finished! Break time.")
		notify("Pomodoro finished!", "Break time.")
		return nil
	},
}

func init() {
	startCmd.Flags().IntVar(&workMinutes, "work", defaultWorkMinutes, "Work duration in minutes")
	startCmd.Flags().IntVar(&breakMinutes, "break", defaultBreakMinutes, "Break duration in minutes")
	startCmd.Flags().StringVar(&taskName, "task", "", "Optional task name")

	rootCmd.AddCommand(startCmd)
}

func printStart(task string) {
	if task == "" {
		fmt.Println("🍅 Pomodoro started")
		return
	}
	fmt.Printf("🍅 Pomodoro started (Task: %s)\n", task)
}

func runCountdown(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := int(duration.Seconds())
	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n⛔️ Pomodoro canceled.")
			return
		case <-ticker.C:
			remaining--
			if remaining < 0 {
				return
			}
			fmt.Printf("⏳ Remaining: %s\n", formatMMSS(remaining))
		}
	}
}

func formatMMSS(totalSeconds int) string {
	if totalSeconds < 0 {
		totalSeconds = 0
	}
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func notify(title, message string) {
	// TODO: Implement OS notifications.
	_ = title
	_ = message
}
