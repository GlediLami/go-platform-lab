# 05: Check 20 components fast

**Ticket:** "The health check runs each component check one after another and takes minutes. Run them in parallel, but never more than N at once (we don't want to hammer the API server)."

**Run:** `make test EX=05` (it uses the race detector)

## Task
- `CheckAll(ctx, names, maxParallel, check)`: run all checks concurrently, at most `maxParallel` at the same time, return every result.
- `Unhealthy(results)`: names with a non-nil error.

## Hints
- `sync.WaitGroup`: `wg.Add(1)` before starting a goroutine, `defer wg.Done()` inside, `wg.Wait()` at the end.
- Semaphore: `sem := make(chan struct{}, maxParallel)`; `sem <- struct{}{}` blocks when full; `<-sem` frees a slot.
- The results map is written by many goroutines: `mu.Lock()` / `mu.Unlock()` around the write.
- Since Go 1.22 each loop iteration has its own variable, so `go func() { ... name ... }()` is safe.

## What you learn
- Goroutines, WaitGroup, Mutex, buffered channels.
- `go test -race` finds data races. If you remove the mutex, run it and read the report.

## In real projects
`golang.org/x/sync/errgroup` (with `SetLimit`) is the library version of this pattern. Operators use it to check or update many components in parallel without overloading the API server.
