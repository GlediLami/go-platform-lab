// Package health checks many things in parallel with a concurrency limit,
// like a platform agent checking the health of all control plane components.
package health

import "context"

// CheckFunc checks one thing. nil means healthy.
type CheckFunc func(ctx context.Context, name string) error

// CheckAll runs check for every name, at most maxParallel at the same time,
// and returns the result for every name (nil = healthy).
//
// TODO(you):
//   - one goroutine per name, wait for all of them (sync.WaitGroup)
//   - limit concurrency with a buffered channel used as a semaphore:
//     sem := make(chan struct{}, maxParallel); sem <- struct{}{} to take, <-sem to give back
//   - many goroutines write to one map: protect it with a sync.Mutex
//   - maxParallel < 1 counts as 1
//
// Run the tests with the race detector: go test -race ./exercises/05-health/
func CheckAll(ctx context.Context, names []string, maxParallel int, check CheckFunc) map[string]error {
	return nil
}

// Unhealthy returns the names whose result is an error.
//
// TODO(you): order doesn't matter (map iteration order is random in Go).
func Unhealthy(results map[string]error) []string {
	return nil
}
