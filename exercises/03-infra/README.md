# 03: Don't create the network twice

**Ticket:** "Our infrastructure code sometimes creates duplicate networks when it runs twice, and errors like `not found` vs `cloud down` look the same in the logs. Make it idempotent and make errors useful."

**Run:** `make test EX=03`

## Task
- `EnsureNetwork`: get first; create only if missing; existing with a different CIDR is an error wrapping `ErrCIDRMismatch`; every other error is wrapped with context.
- `DeleteNetwork`: already gone is success; real errors are returned.

You get `cloud.API` (an interface) and `cloud.Fake` (an in-memory implementation that records calls) in `pkg/cloud`.

## Hints
- `errors.Is(err, cloud.ErrNotFound)` works through wrapping.
- `fmt.Errorf("getting network %q: %w", name, err)`: `%w` keeps the original error inspectable.
- A `switch { case ...: }` without a value reads nicely for "found / not found / other error".

## What you learn
- **Idempotency**: the core property of everything a controller does. It will run again (restarts, retries, resyncs).
- Interfaces + fakes for testing without a cloud.
- Error wrapping and sentinel errors.

## In real projects
Cloud provider integrations (cloud-controller-manager, Cluster API infrastructure providers, cloud-specific operators) wrap the cloud SDK behind interfaces, and their `Reconcile`/`Delete` paths are idempotent exactly like this.
