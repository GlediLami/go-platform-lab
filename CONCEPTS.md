# Concepts: Go for Kubernetes and Gardener work

What you need in your head before (and while) doing the exercises. Each section says which exercise practices it.

---

## 1. Go idioms you'll see everywhere

**Errors are values, and you wrap them.** (03, 04)
```go
if err != nil {
    return fmt.Errorf("creating network %q: %w", name, err)   // %w keeps the original error
}
errors.Is(err, cloud.ErrNotFound)       // true even through several layers of wrapping
```
- Add context ("what was I doing") at each layer. The final message reads like a stack: `ensuring workers Deployment for Cluster garden-dev/dev: Operation cannot be fulfilled...`.
- Sentinel errors (`var ErrNotFound = errors.New(...)`) let callers react to a *kind* of error.
- Kubernetes API errors: `apierrors.IsNotFound(err)`, `client.IgnoreNotFound(err)`.

**Interfaces are small and defined by the user, not the implementer.** (03, 11)
```go
type API interface {                  // only what this controller needs
    GetNetwork(ctx context.Context, name string) (*Network, error)
    CreateNetwork(ctx context.Context, name, cidr string) (*Network, error)
}
```
Real code passes the cloud SDK client; tests pass a fake. That's how Gardener's provider extensions are unit-tested without a cloud.

**`context.Context` is the first parameter of anything that does I/O or can block.** (04, everywhere)
It carries cancellation (manager shutting down), deadlines (timeouts) and the logger. Never store it in a struct; pass it down.

**Concurrency: goroutines + something to wait + something to protect shared data.** (05)
- `sync.WaitGroup` to wait for all goroutines.
- `sync.Mutex` around maps/slices written by several goroutines (Go maps are not safe for concurrent writes).
- A buffered channel as a semaphore to limit parallelism.
- Always run `go test -race` on concurrent code.

**Table-driven tests.** (01, 02, 08)
A slice of cases (`name`, inputs, `want`) and one loop with `t.Run(tt.name, ...)`. Adding a case is one line, and failures name the case.

**Pointers for optional fields.** (08, 10)
Kubernetes uses `*int32` for "may be unset" (e.g. `Deployment.Spec.Replicas`). Use `k8s.io/utils/ptr`: `ptr.To[int32](3)`, `ptr.Deref(p, 1)`.

---

## 2. Kubernetes API objects in Go

Every object has:
- `TypeMeta` (`apiVersion`, `kind`)
- `ObjectMeta` (`name`, `namespace`, `labels`, `annotations`, `ownerReferences`, `finalizers`, `generation`, `resourceVersion`, `deletionTimestamp`, ...)
- `Spec`: what the user **wants** (written by users)
- `Status`: what **is** (written by controllers)

**Scheme** (`pkg/labscheme`): a registry mapping Go types to `group/version/kind`. The client needs it to know that `*v1alpha1.Cluster` is `lab.gardener.cloud/v1alpha1, Kind=Cluster`.

**DeepCopy** (`pkg/apis/lab/v1alpha1/zz_deepcopy.go`): objects from the client's cache are shared, so you must never modify them in place. Generated in real projects (`make generate`), hand-written here so you can read it.

**generation vs resourceVersion:**
- `metadata.generation` increases when the **spec** changes.
- `status.observedGeneration` = the generation the controller last acted on. If they differ, the controller hasn't caught up yet. Gardener uses this everywhere.
- `metadata.resourceVersion` changes on **every** write. It's used for optimistic locking: an update with an old resourceVersion fails with a "conflict", and the controller retries.

**Field validation** (01): `k8s.io/apimachinery/pkg/util/validation/field`. Errors carry a path like `spec.workers[0].maximum`, which is exactly what `kubectl` prints when Gardener rejects a Shoot.

---

## 3. The controller-runtime client (06, 07)

```go
c.Get(ctx, client.ObjectKey{Namespace: ns, Name: name}, obj)
c.List(ctx, &corev1.PodList{}, client.InNamespace(ns), client.MatchingLabels{"app": "x"})
c.Create(ctx, obj)
c.Update(ctx, obj)                            // whole object, needs current resourceVersion
c.Patch(ctx, obj, client.MergeFrom(original)) // only the diff: fewer conflicts
c.Delete(ctx, obj)
c.Status().Update(ctx, obj)                   // status subresource
```

- In a manager, `mgr.GetClient()` **reads from a cache** (watched by informers) and **writes to the API server**. Reads are cheap, but can be slightly stale.
- `controllerutil.CreateOrUpdate(ctx, c, obj, mutate)`: get, run `mutate` to set the desired fields, then create or update only if something changed. Returns `Created`, `Updated` or `None`. Set **only the fields you own** inside `mutate`.
- Tests use `sigs.k8s.io/controller-runtime/pkg/client/fake`: an in-memory client with the same interface.

---

## 4. Anatomy of a controller (08 to 12)

```go
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
    obj := &v1alpha1.Cluster{}
    if err := r.Client.Get(ctx, req.NamespacedName, obj); err != nil {
        return reconcile.Result{}, client.IgnoreNotFound(err)   // deleted meanwhile: done
    }
    if !obj.DeletionTimestamp.IsZero() {
        return r.delete(ctx, obj)                               // cleanup path (finalizer)
    }
    // read desired (spec) + actual (cluster/cloud), fix the difference, write status
    return reconcile.Result{}, nil
}
```

**Level-based, not event-based.** `req` only has namespace/name, not "what changed". You always look at the full current state and converge. So Reconcile must be **idempotent**: running it 1 or 100 times gives the same result.

**Return values:**
- `return reconcile.Result{}, nil`: done, wait for the next change.
- `return reconcile.Result{}, err`: retry with exponential backoff. Only for things a retry can fix (API errors, cloud hiccups).
- `return reconcile.Result{RequeueAfter: time.Minute}, nil`: check again later (polling something external).
- User mistakes (invalid spec) are **not** errors: report them on status, don't spin.

**Wiring** (`SetupWithManager`):
```go
ctrl.NewControllerManagedBy(mgr).
    For(&v1alpha1.Cluster{}).        // reconcile when a Cluster changes
    Owns(&appsv1.Deployment{}).      // ...or a Deployment the Cluster controls changes
    Complete(r)
```

**Owner references** (10): `controllerutil.SetControllerReference(owner, obj, scheme)`. Two effects: deleting the owner garbage-collects the object, and `Owns()` maps events on the object back to the owner. That's how drift ("someone scaled my Deployment") gets fixed automatically.

**Finalizers** (11): a string in `metadata.finalizers`. While present, `kubectl delete` only sets `deletionTimestamp`, and the object stays. The controller cleans up external things (cloud resources), then removes its finalizer, then Kubernetes deletes the object. Add the finalizer **before** creating anything external.

**Status and conditions** (09): `[]metav1.Condition` with `type`, `status` (True/False/Unknown), `reason` (CamelCase), `message`, `lastTransitionTime`, `observedGeneration`. Use `meta.SetStatusCondition` / `meta.FindStatusCondition`. Gardener's Shoot shows conditions like `APIServerAvailable`, `ControlPlaneHealthy`, `EveryNodeReady`.

**Manager** (12): one process, one shared cache and client, many controllers, leader election, health probes, graceful shutdown. `cmd/gardenlet` is a manager with dozens of controllers.

---

## 5. Testing levels (like Gardener)

| Level | What | Tools | Here |
|-------|------|-------|------|
| Unit | one function or reconciler, fake dependencies | `testing`, fake client, fakes/mocks | 01 to 11 |
| Integration | real kube-apiserver + etcd, no nodes | envtest (`test/integration` in Gardener) | 12 |
| End-to-end | a whole Gardener landscape | kind + local setup (`test/e2e`) | the Gardener lab repo |

Gardener writes tests with **Ginkgo** (`Describe`, `Context`, `It`, `BeforeEach`, `JustBeforeEach`) and **Gomega** (`Expect(x).To(Equal(y))`, `Eventually(...).Should(...)`), and mocks interfaces with **gomock** (`go.uber.org/mock`). Exercises 09 to 12 use Ginkgo/Gomega so the style feels familiar when you open Gardener's `_test.go` files.

`Eventually` matters for integration tests: controllers are asynchronous, so you poll until the expectation holds or a timeout hits.

---

## 6. How Gardener code is organized (where to look)

- `pkg/apis/core/v1beta1/`: API types (`types_shoot.go`), generated deepcopy
- `pkg/api/core/validation/`: validation (`shoot.go`), the big brother of exercise 01
- `pkg/gardenlet/controller/`: gardenlet controllers (shoot, care, ...)
- `pkg/controllerutils/`: helpers for patching, finalizers, reconcilers
- `pkg/utils/flow/`: run tasks as a dependency graph, in parallel
- `pkg/utils/retry/`: retry/wait helpers
- `pkg/utils/kubernetes/health/`: health checks for Deployments, Pods, ...
- `extensions/pkg/controller/`: the extension library (Infrastructure, Worker, ControlPlane actuators)
- `pkg/features/`: feature gates
- `cmd/`: the binaries (gardenlet, gardener-operator, ...)
- `test/integration/`, `test/e2e/`: integration and end-to-end tests

---

## 7. How a feature lands in Gardener (the work itself)

1. **Discuss**: issue or proposal (GEP for big changes). Agree on the API.
2. **API**: add fields to types, `make generate` (deepcopy, CRDs, docs), validation + defaulting.
3. **Feature gate** for risky changes (`pkg/features`, documented in `docs/deployment/feature_gates.md`): alpha (off), beta (on), GA.
4. **Controller logic**: reconcile + delete paths, idempotent, status/conditions.
5. **Tests**: unit (Ginkgo, fake client, gomock), integration (envtest), e2e if user-visible.
6. **Docs** in `docs/`, a **release note** in the PR template.
7. `make verify` (check, format, test, test-integration) green, then review.

Smaller day-to-day work: fixing a bug in a reconciler, adding a field to a provider extension, improving a health check, adding a condition. All of it is the pattern from exercises 08 to 12.
