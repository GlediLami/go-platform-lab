# 12: Ship it: one manager

**Ticket:** "We have four controllers in four packages. Wire them into one binary, prove they work together against a real API server, then run them on a cluster."

**Run:** `make test EX=12` (first run downloads `kube-apiserver` + `etcd` for envtest)

## Task
Implement `Setup(mgr, cloudAPI)` in `manager.go`: register **your** reconcilers from exercises 08, 09 (with your `ValidateCluster` from 01), 10 and 11.

`cmd/main.go` is already written: it creates the manager and calls `Setup`.

## What the test does
`manager_test.go` is an **integration test**. envtest starts a real `kube-apiserver` and `etcd`, installs `config/crd/`, starts your manager, and checks with `Eventually`:
1. `SpecValid=True` on a good Cluster
2. `dev-workers` Deployment with 3 replicas
3. network created + finalizer added
4. scale the Deployment to 10 by hand -> back to 3 (`Owns` in action)
5. hibernate -> 0 replicas
6. delete -> network deleted, Cluster gone
7. a bad Cluster gets `SpecValid=False` with the field path in the message

If this passes, your controllers work together for real.

## Then run it on kind
See section 5 of the main README: `kind create cluster --name lab`, apply the CRD, `go run ./exercises/12-manager/cmd`, apply `config/samples/cluster.yaml` and watch the logs while you scale, hibernate and delete.

## What you learn
- The Manager: shared cache, many controllers, graceful shutdown.
- Integration tests with envtest: real API semantics (status subresource, finalizers, resourceVersion conflicts) without a full cluster. Note what's missing: no kube-controller-manager, so no garbage collection and no pods.
- Two controllers writing the same object cause **conflicts**: look for `the object has been modified` in the logs. controller-runtime retries them, and your idempotent code makes that safe.

## In Gardener
`cmd/gardenlet` builds a manager and adds all gardenlet controllers. `test/integration/` has envtest suites for most controllers, set up just like `suite_test.go` here.
