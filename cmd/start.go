package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"pomodoro/internal/notify"
	"pomodoro/internal/state"
	"pomodoro/internal/timer"

	"github.com/spf13/cobra"
)

const (
	defaultWorkMinutes     = 25
	defaultBreakMinutes    = 5
	defaultLongBreakMinute = 15
	defaultSessionsTotal   = 4
	stateVersion           = 1
	persistEverySeconds    = 5
)

var (
	workMinutes     int
	breakMinutes    int
	longBreakMinute int
	sessionsTotal   int
	taskName        string
	forceStart      bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a pomodoro work session",
	RunE: func(cmd *cobra.Command, args []string) error {
		if workMinutes <= 0 {
			return fmt.Errorf("work minutes must be greater than 0")
		}
		if breakMinutes <= 0 {
			return fmt.Errorf("break minutes must be greater than 0")
		}
		if longBreakMinute <= 0 {
			return fmt.Errorf("long-break minutes must be greater than 0")
		}
		if sessionsTotal <= 0 {
			return fmt.Errorf("sessions must be greater than 0")
		}

		existing, err := state.Load()
		if err != nil {
			return fmt.Errorf("load state: %w", err)
		}
		if existing != nil && existing.Running && !forceStart {
			return fmt.Errorf("A session is already running. Use --force to restart.")
		}

		workSeconds := workMinutes * 60
		breakSeconds := breakMinutes * 60
		longBreakSeconds := longBreakMinute * 60

		now := time.Now().UTC()
		startedAt := now.Format(time.RFC3339)
		endsAt := now.Add(time.Duration(workSeconds) * time.Second).Format(time.RFC3339)

		current := &state.State{
			Version:          stateVersion,
			Running:          true,
			Mode:             "work",
			Task:             taskName,
			SessionIndex:     1,
			SessionsTotal:    sessionsTotal,
			WorkSeconds:      workSeconds,
			BreakSeconds:     breakSeconds,
			LongBreakSeconds: longBreakSeconds,
			StartedAt:        startedAt,
			EndsAt:           endsAt,
			Paused:           false,
			PausedAt:         nil,
			RemainingSeconds: workSeconds,
		}

		if err := state.Save(current); err != nil {
			return fmt.Errorf("save state: %w", err)
		}

		printStart(taskName)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
		defer stop()

		onTick := func(remaining int) {
			current.RemainingSeconds = remaining
			fmt.Printf("⏳ Remaining: %s\n", formatMMSS(remaining))
		}
		onPersist := func(remaining int) error {
			current.RemainingSeconds = remaining
			return state.Save(current)
		}

		remaining, completed, err := timer.RunCountdown(
			ctx,
			workSeconds,
			onTick,
			persistEverySeconds,
			onPersist,
		)
		if err != nil && err != context.Canceled {
			return fmt.Errorf("countdown: %w", err)
		}

		if err == context.Canceled {
			current.Running = false
			current.RemainingSeconds = remaining
			if saveErr := state.Save(current); saveErr != nil {
				return fmt.Errorf("save canceled state: %w", saveErr)
			}
			fmt.Println("\n⛔️ Pomodoro canceled.")
			return nil
		}

		if completed {
			fmt.Println("✅ Pomodoro finished! Break time.")
			notify.Notify("Pomodoro finished!", "Break time.")

			nextMode := "break"
			nextSeconds := breakSeconds
			if current.SessionIndex%sessionsTotal == 0 {
				nextMode = "long_break"
				nextSeconds = longBreakSeconds
			}

			nextStart := time.Now().UTC()
			current.Mode = nextMode
			current.Running = false
			current.StartedAt = nextStart.Format(time.RFC3339)
			current.EndsAt = nextStart.Add(time.Duration(nextSeconds) * time.Second).Format(time.RFC3339)
			current.RemainingSeconds = nextSeconds

			if err := state.Save(current); err != nil {
				return fmt.Errorf("save break state: %w", err)
			}
		}

		return nil
	},
}

func init() {
	startCmd.Flags().IntVar(&workMinutes, "work", defaultWorkMinutes, "Work duration in minutes")
	startCmd.Flags().IntVar(&breakMinutes, "break", defaultBreakMinutes, "Break duration in minutes")
	startCmd.Flags().IntVar(&longBreakMinute, "long-break", defaultLongBreakMinute, "Long break duration in minutes")
	startCmd.Flags().IntVar(&sessionsTotal, "sessions", defaultSessionsTotal, "Number of work sessions before a long break")
	startCmd.Flags().StringVar(&taskName, "task", "", "Optional task name")
	startCmd.Flags().BoolVar(&forceStart, "force", false, "Restart even if a session is running")

	rootCmd.AddCommand(startCmd)
}

func printStart(task string) {
	if task == "" {
		fmt.Println("🍅 Pomodoro started")
		return
	}
	fmt.Printf("🍅 Pomodoro started (Task: %s)\n", task)
}

func formatMMSS(totalSeconds int) string {
	if totalSeconds < 0 {
		totalSeconds = 0
	}
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
