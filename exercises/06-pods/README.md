# 06: Which pods are broken?

**Ticket:** "On-call keeps typing `kubectl get pods | grep -v Running`. Give us a function that returns the unhealthy pods of a namespace."

**Run:** `make test EX=06`

## Task
`UnhealthyPods(ctx, c, namespace)` returns the **sorted** names of pods that are not healthy:
- `Succeeded` is healthy
- `Running` is healthy only if every container is ready
- everything else (`Pending`, `Failed`, `Unknown`) is unhealthy

Pods in other namespaces must be ignored.

## Hints
```go
podList := &corev1.PodList{}
err := c.List(ctx, podList, client.InNamespace(namespace))
for _, pod := range podList.Items { ... }
```
- `corev1` is `k8s.io/api/core/v1`.
- `pod.Status.Phase`, `pod.Status.ContainerStatuses[i].Ready`.

## What you learn
- The controller-runtime `client.Client` interface: the same in production and tests.
- The fake client (`fake.NewClientBuilder().WithObjects(...).Build()`): Kubernetes without a cluster.
- Go types for Kubernetes objects.

## In real projects
Readiness logic like this sits behind `kubectl rollout status`, operator health checks and status conditions such as `Available`.
