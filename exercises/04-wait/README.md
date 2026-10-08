# 04: `kubectl wait`, in Go

**Ticket:** "We need a helper that waits until something is ready (a cluster, a condition), checks right away, polls at an interval, stops on real errors and gives up when the context times out."

**Run:** `make test EX=04`

## Task
Implement `Until(ctx, interval, cond)`:
- check immediately, then every `interval`
- `true` -> return nil
- error -> return it (wrapped), stop
- `ctx` done -> return an error that wraps **both** `ErrTimeout` and `ctx.Err()`

## Hints
```go
ticker := time.NewTicker(interval)
defer ticker.Stop()
for {
    // check...
    select {
    case <-ctx.Done(): // give up
    case <-ticker.C:   // try again
    }
}
```
- Go 1.20+ supports two `%w` in one `fmt.Errorf`.
- Callers set the timeout: `ctx, cancel := context.WithTimeout(ctx, 15*time.Minute)`.

## What you learn
- `context.Context` for cancellation and deadlines.
- `select` over channels, tickers.
- It's what a `wait-for` script or `kubectl wait --for=condition=...` does, as a Go function.

## In real projects
`kubectl wait`, `wait.PollUntilContextTimeout` in `k8s.io/apimachinery/pkg/util/wait`, and Gomega's `Eventually` in tests all do this job.
