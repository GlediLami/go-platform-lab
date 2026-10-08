# 09: Tell the user what's wrong

**Ticket:** "When a Cluster spec is invalid, nothing happens and users don't know why. Report it on the Cluster status as a condition, like Shoots do."

**Run:** `make test EX=09`

From here on the tests use **Ginkgo/Gomega**, like Gardener. Open `status_test.go` and read it first: `Describe`, `BeforeEach`, `JustBeforeEach`, `It`, `Expect(...).To(...)`.

## Task
Implement `Reconcile`:
- Validation errors -> condition `SpecValid` = `False`, reason `InvalidSpec`, message = the errors
- No errors -> `SpecValid` = `True`, reason `Valid`
- Set `ObservedGeneration` on the condition and on `status`
- Write with the **status** client
- An invalid spec is **not** a returned error

Validation is injected (`r.Validate`), so this test doesn't depend on your exercise 01. In exercise 12, you'll plug your `ValidateCluster` in.

## Hints
- `meta.SetStatusCondition(&cluster.Status.Conditions, condition)` (`k8s.io/apimachinery/pkg/api/meta`): adds or replaces by `Type`, sets `LastTransitionTime` when the status flips.
- `r.Client.Status().Update(ctx, cluster)`: a plain `Update` ignores status, because status is a subresource.
- `errs.ToAggregate().Error()` turns a `field.ErrorList` into one message.

## What you learn
- Spec is written by users, status by controllers.
- Conditions: the standard way to report state (`kubectl describe` shows them).
- Dependency injection via a function field.
- Ginkgo structure.

## In Gardener
The care controller (`pkg/gardenlet/controller/shoot/care`) sets Shoot conditions like `APIServerAvailable` and `EveryNodeReady`: the ones `wait-for.sh` waited for in your setup.
