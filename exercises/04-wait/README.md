# 04: wait-for.sh, in Go

**Ticket:** "We need a helper that waits until something is ready (a Seed, a Shoot condition), checks right away, polls at an interval, stops on real errors and gives up when the context times out."

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
- This is `hack/usage/wait-for.sh` (from your Gardener setup) as a Go function.

## In Gardener
`pkg/utils/retry` and `k8s.io/apimachinery/pkg/util/wait` (`wait.PollUntilContextTimeout`). In tests, Gomega's `Eventually` does the same job.
