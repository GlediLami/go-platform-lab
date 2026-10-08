# Exercises

Each folder is a ticket. Read its `README.md`, run `make test EX=NN`, write code until it's green.

| # | Ticket | You learn |
|---|--------|-----------|
| 01 | [Reject broken Clusters](01-validation/README.md) | structs, `field.ErrorList`, field paths |
| 02 | [How many machines?](02-replicas/README.md) | table-driven tests, maps, zero values |
| 03 | [Don't create the network twice](03-infra/README.md) | interfaces, fakes, error wrapping, idempotency |
| 04 | [`kubectl wait`, in Go](04-wait/README.md) | `context`, timeouts, `select`, tickers |
| 05 | [Check 20 components fast](05-health/README.md) | goroutines, `WaitGroup`, mutex, semaphore, `-race` |
| 06 | [Which pods are broken?](06-pods/README.md) | controller-runtime client, fake client |
| 07 | [Create or update, safely](07-configmap/README.md) | `CreateOrUpdate`, owning only your fields |
| 08 | [Your first reconciler](08-reconciler/README.md) | `Reconcile`, not-found, patch, user vs system errors |
| 09 | [Tell the user what's wrong](09-status/README.md) | status subresource, conditions, Ginkgo |
| 10 | [The Cluster owns its workers](10-owned/README.md) | owner references, `Owns()`, hibernation |
| 11 | [Clean up before you go](11-finalizer/README.md) | finalizers, deletion flow |
| 12 | [Ship it: one manager](12-manager/README.md) | manager wiring, envtest integration test, kind |

Rules:
- Don't change the `_test.go` files (except where a test says `TODO(you)`).
- Read the failing test before writing code: tests are the spec.
- Stuck for 20 minutes? Look at the hint in the ticket, then at `solutions/NN-*/`.
