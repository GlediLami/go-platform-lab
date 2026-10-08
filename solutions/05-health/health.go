// Package health checks many things in parallel with a concurrency limit,
// like a platform agent checking the health of all control plane components.
package health

import (
	"context"
	"sync"
)

// CheckFunc checks one thing. nil means healthy.
type CheckFunc func(ctx context.Context, name string) error

// CheckAll runs check for every name, at most maxParallel at the same time,
// and returns the result for every name (nil = healthy).
func CheckAll(ctx context.Context, names []string, maxParallel int, check CheckFunc) map[string]error {
	if maxParallel < 1 {
		maxParallel = 1
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results = make(map[string]error, len(names))
		sem     = make(chan struct{}, maxParallel) // a semaphore: at most maxParallel tokens
	)

	for _, name := range names {
		wg.Add(1)
		go func() {
			defer wg.Done()

			sem <- struct{}{}        // take a token (blocks while maxParallel are running)
			defer func() { <-sem }() // give it back

			err := check(ctx, name)

			mu.Lock()
			results[name] = err
			mu.Unlock()
		}()
	}

	wg.Wait()
	return results
}

// Unhealthy returns the names whose result is an error.
func Unhealthy(results map[string]error) []string {
	var out []string
	for name, err := range results {
		if err != nil {
			out = append(out, name)
		}
	}
	return out
}
