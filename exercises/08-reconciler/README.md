# 08: Your first reconciler

**Ticket:** "Teams want to guarantee a minimum number of replicas on some Deployments, even if someone scales them down by hand. Add an annotation `lab.gardener.cloud/min-replicas: \"3\"` and enforce it."

**Run:** `make test EX=08`

## Task
Implement `Reconcile`:
1. Get the Deployment. Not found -> done, **no error**.
2. No annotation -> done.
3. Invalid annotation (not a number, negative) -> log and done, **no error**.
4. Replicas below minimum -> patch them up. (`nil` replicas means 1.)
5. Return an error only for real API failures.

`SetupWithManager` is already written: read it.

## Hints
- `client.IgnoreNotFound(err)`
- `strconv.ParseInt(value, 10, 32)`
- `ptr.Deref(deployment.Spec.Replicas, 1)`, `ptr.To(int32(n))` from `k8s.io/utils/ptr`
- Patch only what changed:
  ```go
  patch := client.MergeFrom(deployment.DeepCopy())
  deployment.Spec.Replicas = ptr.To(int32(minReplicas))
  err := r.Client.Patch(ctx, deployment, patch)
  ```
- Logger: `log := logf.FromContext(ctx)` (`sigs.k8s.io/controller-runtime/pkg/log`), then `log.Info("msg", "key", value)`.

## What you learn
- The shape of every controller: get, compare, fix, return.
- **System errors vs user errors**: return errors only when a retry can help. A typo in an annotation won't fix itself; retrying it forever just spams logs.
- Level-based thinking: you don't know *what* changed, you just converge.

## In Gardener
Every controller under `pkg/gardenlet/controller/` starts like this. `pkg/controllerutils/reconciler/` has shared helpers.
