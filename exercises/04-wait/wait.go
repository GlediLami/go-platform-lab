// Package wait polls a condition until it is true, like hack/usage/wait-for.sh
// in Gardener or wait.PollUntilContextTimeout in k8s.io/apimachinery.
package wait

import (
	"context"
	"errors"
	"time"
)

// ErrTimeout is returned when the context ends before the condition is true.
var ErrTimeout = errors.New("timed out waiting for the condition")

// ConditionFunc reports whether the wait is over. Returning an error stops the wait.
type ConditionFunc func(ctx context.Context) (done bool, err error)

// Until checks cond immediately and then every interval until it returns true,
// returns an error, or ctx is done.
//
// TODO(you):
//   - call cond(ctx) right away (no waiting before the first check)
//   - done        -> return nil
//   - error       -> return it wrapped (stop immediately)
//   - otherwise wait for the next tick OR ctx.Done(), whichever comes first
//     (use time.NewTicker and select)
//   - ctx done    -> return an error wrapping BOTH ErrTimeout and ctx.Err():
//     fmt.Errorf("%w: %w", ErrTimeout, ctx.Err())
func Until(ctx context.Context, interval time.Duration, cond ConditionFunc) error {
	return errors.New("not implemented")
}
