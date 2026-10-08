# 10: The Cluster owns its workers

**Ticket:** "For each Cluster, run a Deployment `<name>-workers` with one pod per worker machine (sum of pool minimums, 0 when hibernated). If anyone scales or edits it, put it back. Deleting the Cluster must delete it."

**Run:** `make test EX=10`

## Task
- `DesiredReplicas(cluster)`: sum of `Minimum` over pools, 0 if hibernated.
- `Reconcile`: `CreateOrUpdate` the Deployment, make the Cluster its **controller owner**, write `status.replicas`.

`SetupWithManager` uses `Owns(&appsv1.Deployment{})`: read why in the comment.

## Hints
- The selector of a Deployment is **immutable**: set it only when it's nil.
- Selector `matchLabels` must match the pod template labels.
- `controllerutil.SetControllerReference(cluster, deployment, r.Scheme)` as the last line of the mutate func.
- Update status only if the value changed (avoids pointless writes).

## What you learn
- Owner references: garbage collection + event mapping (`Owns`).
- Building a full child object in Go (selector, template, containers).
- Hibernation as "desired replicas 0".

## In real projects
A Deployment owns ReplicaSets, which own Pods: same mechanism. Cluster API's MachineDeployment owns MachineSets, which own Machines.
