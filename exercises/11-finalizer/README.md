# 11: Clean up before you go

**Ticket:** "When users delete a Cluster, its cloud network stays behind and costs money. Create a network per Cluster and make sure it's deleted before the Cluster disappears."

**Run:** `make test EX=11`

## Task
- `NetworkName`: `cluster--<namespace>--<name>`
- Normal reconcile: add finalizer `lab.gardener.cloud/network` **first**, then ensure the network exists.
- Deletion (`deletionTimestamp` set): delete the network (already gone = fine), then remove the finalizer. If the cloud fails, return the error and **keep** the finalizer.

## Hints
- `controllerutil.ContainsFinalizer`, `AddFinalizer`, `RemoveFinalizer`
- Patch: `patch := client.MergeFrom(cluster.DeepCopy())`, change, `r.Client.Patch(ctx, cluster, patch)`
- Split `Reconcile` into `reconcile()` and `delete()` helpers: it reads much better.
- Why finalizer first: if you create the network first and the process crashes before adding the finalizer, a delete would leave the network orphaned forever.

## What you learn
- The deletion flow: `kubectl delete` -> `deletionTimestamp` -> your cleanup -> finalizer removed -> object gone.
- Why "stuck in Terminating" happens (a finalizer nobody removes).
- Ordering matters for crash safety.

## In Gardener
Shoot deletion runs a long flow (delete workers, control plane, infrastructure, DNS...) before removing the `gardener` finalizer. Helpers: `pkg/controllerutils/finalizers.go`. In the Gardener lab, scenario 12 showed the same thing from the outside.
