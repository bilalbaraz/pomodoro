package timer

import (
	"context"
	"fmt"
	"time"
)

// RunCountdown ticks every second until remaining reaches zero or the context is canceled.
// It returns the remaining seconds, a completed flag, and any error encountered.
func RunCountdown(
	ctx context.Context,
	totalSeconds int,
	onTick func(remaining int),
	persistEvery int,
	onPersist func(remaining int) error,
) (int, bool, error) {
	if totalSeconds < 0 {
		return 0, false, fmt.Errorf("totalSeconds must be >= 0")
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	remaining := totalSeconds
	ticks := 0

	for {
		select {
		case <-ctx.Done():
			return remaining, false, ctx.Err()
		case <-ticker.C:
			remaining--
			if remaining < 0 {
				return 0, true, nil
			}

			if onTick != nil {
				onTick(remaining)
			}

			ticks++
			if persistEvery > 0 && ticks%persistEvery == 0 && onPersist != nil {
				if err := onPersist(remaining); err != nil {
					return remaining, false, err
				}
			}

			if remaining == 0 {
				return 0, true, nil
			}
		}
	}
}
