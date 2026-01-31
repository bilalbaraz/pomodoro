package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State represents the persisted pomodoro session state.
type State struct {
	Version          int     `json:"version"`
	Running          bool    `json:"running"`
	Mode             string  `json:"mode"`
	Task             string  `json:"task"`
	SessionIndex     int     `json:"session_index"`
	SessionsTotal    int     `json:"sessions_total"`
	WorkSeconds      int     `json:"work_seconds"`
	BreakSeconds     int     `json:"break_seconds"`
	LongBreakSeconds int     `json:"long_break_seconds"`
	StartedAt        string  `json:"started_at"`
	EndsAt           string  `json:"ends_at"`
	Paused           bool    `json:"paused"`
	PausedAt         *string `json:"paused_at"`
	RemainingSeconds int     `json:"remaining_seconds"`
}

// ResolveStatePath returns the full path to the state.json file.
func ResolveStatePath() (string, error) {
	xdg := os.Getenv("XDG_STATE_HOME")
	if xdg != "" {
		return filepath.Join(xdg, "pomodoro", "state.json"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".pomodoro", "state.json"), nil
}

// Load reads the state from disk. Returns (nil, nil) if not found.
func Load() (*State, error) {
	path, err := ResolveStatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read state: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	return &s, nil
}

// Save writes the state to disk, creating directories if needed.
func Save(s *State) error {
	path, err := ResolveStatePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	return nil
}
