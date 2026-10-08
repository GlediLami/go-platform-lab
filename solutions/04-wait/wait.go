// Package wait polls a condition until it is true, like `kubectl wait` or
// wait.PollUntilContextTimeout in k8s.io/apimachinery.
package wait

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrTimeout is returned when the context ends before the condition is true.
var ErrTimeout = errors.New("timed out waiting for the condition")

// ConditionFunc reports whether the wait is over. Returning an error stops the wait.
type ConditionFunc func(ctx context.Context) (done bool, err error)

// Until checks cond immediately and then every interval until it returns true,
// returns an error, or ctx is done.
func Until(ctx context.Context, interval time.Duration, cond ConditionFunc) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		done, err := cond(ctx)
		if err != nil {
			return fmt.Errorf("condition failed: %w", err)
		}
		if done {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrTimeout, ctx.Err())
		case <-ticker.C:
		}
	}
}
