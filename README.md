# 🍅 pomodoro — Focus sessions from the terminal

A small Go CLI that runs Pomodoro work/break sessions, tracks state on disk, and supports pause/resume.

> State is stored in `~/.pomodoro/state.json` (or `$XDG_STATE_HOME/pomodoro/state.json`).

## Installation (Homebrew)
```bash
brew tap bilalbaraz/tap
brew install pomodoro
```

## Quickstart
```bash

# start a default work session (25 min)
pomodoro start

# start with a task and custom durations
pomodoro start --task "Auth Refactor" --work 50 --break 10

# pause and resume
pomodoro pause
pomodoro resume

# check status
pomodoro status
```

## Command Surface
- **Start:** `pomodoro start`
- **Pause:** `pomodoro pause`
- **Resume:** `pomodoro resume`
- **Stop:** `pomodoro stop`
- **Status:** `pomodoro status`
- **Stats:** `pomodoro stats`

## Flags
- **Start:** `--work`, `--break`, `--long-break`, `--sessions`, `--task`, `--force`
- **Status:** `--json`

## Configuration
No config file. All behavior is controlled via CLI flags and the state file.

## Notes
- Errors are printed to stderr and exit with code 1 on invalid input.
