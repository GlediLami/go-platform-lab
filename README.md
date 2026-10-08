# Go Gardener Lab

Learn to write the Go code a Gardener / SKE team writes: validation, idempotent cloud calls, retries, concurrency, and **Kubernetes controllers**, step by step.

Same idea as the Gardener lab: every exercise is a small **ticket**. Here you don't edit YAML, you write Go until the tests go green.

- New to the ideas? Read [CONCEPTS.md](CONCEPTS.md) first: Go idioms used at Gardener, how a controller works, how Gardener's code is organized.
- Stuck on an error? See [TROUBLESHOOTING.md](TROUBLESHOOTING.md).
- The exercises: [exercises/README.md](exercises/README.md).

---

## 1. Setup on macOS (once)

```bash
brew install go
go version                       # 1.26 or newer
```

Editor: VS Code with the **Go** extension (it installs `gopls`, the language server). It gives you autocomplete, "go to definition" into the Kubernetes libraries, and format on save.

Get the repo:
```bash
cd ~/dev
git clone https://github.com/GlediLami/go-gardener-lab.git
cd go-gardener-lab
go mod download
```

Check everything works (runs the reference solutions, all should be `ok`):
```bash
make solutions
```
The first run downloads Go modules and, for exercise 12, a real `kube-apiserver` + `etcd` (envtest). That takes a minute.

---

## 2. How an exercise works

```
exercises/01-validation/
├── README.md            the ticket: what, why, hints, what you learn
├── validation.go        YOUR code: stubs with TODO(you) comments
└── validation_test.go   the tests (don't change them, except where a TODO says so)

solutions/01-validation/ the reference implementation with the same tests
```

1. Read `exercises/NN-*/README.md`.
2. Run the tests and watch them fail:
   ```bash
   make test EX=01
   ```
3. Write code in the `.go` file until they pass. Read the failure messages: they tell you what's expected.
4. Compare with `solutions/NN-*/` afterwards (or when stuck for more than ~20 minutes):
   ```bash
   diff exercises/01-validation/validation.go solutions/01-validation/validation.go
   ```
5. Commit your progress:
   ```bash
   git add -A && git commit -m "exercise 01 done" && git push
   ```

---

## 3. Commands you'll use

```bash
make test EX=03            # your exercise 03 (with the race detector)
make solution EX=03        # the reference solution's tests
make test-all              # all your exercises
make check                 # gofmt + go vet + tests in sync (like `make check` in Gardener)

go test ./exercises/03-infra/ -run TestEnsureNetwork_IsIdempotent -v   # one test, verbose
go test -race ./exercises/05-health/                                   # race detector
go test -cover ./exercises/01-validation/                               # coverage %
go doc sigs.k8s.io/controller-runtime/pkg/controller/controllerutil CreateOrUpdate   # read API docs
gofmt -w .                 # format (VS Code does this on save)
```

Ginkgo tests (exercises 09 to 12) also run with plain `go test`. To focus one spec, put an `F` in front: `FIt(...)`, then remove it again.

---

## 4. The path

| # | Exercise | Skill | Like in Gardener |
|---|----------|-------|------------------|
| 01 | Validation | structs, error lists, field paths | `pkg/api/core/validation/shoot.go` |
| 02 | Replicas | table-driven tests | worker pool min/max |
| 03 | Idempotent infra | interfaces, fakes, error wrapping | provider extension Infrastructure actuator |
| 04 | Wait | context, timeouts, select | `pkg/utils/retry`, `hack/usage/wait-for.sh` |
| 05 | Health | goroutines, mutex, semaphore | `pkg/utils/flow`, health checks |
| 06 | Pods | controller-runtime client, fake client | `pkg/utils/kubernetes/health` |
| 07 | ConfigMap | CreateOrUpdate, owning only your fields | `pkg/controllerutils/patch.go` |
| 08 | First reconciler | Reconcile, not-found, patch | every controller |
| 09 | Status | conditions, status subresource, Ginkgo | Shoot conditions (`care` controller) |
| 10 | Owned objects | owner references, `Owns()` | Worker -> MachineDeployment |
| 11 | Finalizer | cleanup before deletion | Shoot deletion / infrastructure cleanup |
| 12 | Manager | wire it all, integration test (envtest), run on kind | `cmd/gardenlet`, `test/integration` |

01 to 05 are plain Go. 06 to 12 are Kubernetes. Do them in order: each one uses what the previous taught.

---

## 5. Exercise 12 on a real cluster (optional)

After 12 passes, run your controllers against a real kind cluster:
```bash
brew install kind
kind create cluster --name lab
kubectl apply -f config/crd/
go run ./exercises/12-manager/cmd          # keep this tab open: controller logs
```
In a second tab:
```bash
kubectl apply -f config/samples/cluster.yaml
kubectl get clusters,deployments            # Cluster dev, Deployment dev-workers with 3 replicas
kubectl scale deployment dev-workers --replicas=10 && kubectl get deploy dev-workers -w   # watch it go back to 3
kubectl patch cluster dev --type merge -p '{"spec":{"hibernated":true}}'
kubectl delete cluster dev                  # watch the finalizer log line in tab 1
kind delete cluster --name lab
```
Don't run this against the Gardener kind cluster (`gardener-local`): use the separate `lab` cluster.

---

## 6. After this lab

- Read real Gardener code for each pattern (paths in the table above and in each ticket).
- `website-operator` (kubebuilder) shows the same ideas with generated code (CRDs from Go types, `make manifests`).
- Next level: webhooks (validation/defaulting on admission), `gomock` mocks (`go.uber.org/mock`, used across Gardener), and Gardener's extension library (`extensions/pkg/controller`).
