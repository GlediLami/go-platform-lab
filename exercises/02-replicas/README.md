# 02: How many machines?

**Ticket:** "The worker controller needs one function that says how many machines a pool should run, respecting min/max and hibernation."

**Run:** `make test EX=02`

## Task
1. `DesiredReplicas(pool, current, hibernated)`: 0 when hibernated, otherwise `current` clamped into `[Minimum, Maximum]`.
2. `TotalDesired(spec, current)`: sum over all pools. A pool missing from the `current` map has 0 machines.
3. In `replicas_test.go`, add **two more test cases** of your own where the `TODO(you)` is.

## Hints
- Reading a missing key from a map returns the zero value: `current["nope"] == 0`. No special case needed.
- Reuse `DesiredReplicas` inside `TotalDesired`.

## What you learn
- Table-driven tests: the standard Go test style. Each case is one line; `t.Run` names the failing case.
- Writing your own test cases: think about edges (0, equal to min, equal to max).

## In Gardener
Worker pools have `minimum`/`maximum` (cluster-autoscaler bounds), and hibernation scales workers to zero. See `Worker` in `pkg/apis/core/v1beta1/types_shoot.go`.
