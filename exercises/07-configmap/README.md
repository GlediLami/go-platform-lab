# 07: Create or update, safely

**Ticket:** "Our controller overwrites labels other teams put on its ConfigMap, and it writes to the API server on every run even when nothing changed. Fix both."

**Run:** `make test EX=07`

## Task
`EnsureConfigMap(ctx, c, namespace, name, data)`:
- the ConfigMap exists with exactly `data`
- it has the label `app.kubernetes.io/managed-by=go-gardener-lab`
- labels set by others survive
- return `Created`, `Updated` or `None` (nothing changed)

## Hints
```go
cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}
result, err := controllerutil.CreateOrUpdate(ctx, c, cm, func() error {
    // set the desired state on cm here
    return nil
})
```
- Inside the func, `cm` is either empty (create) or the current object (update).
- Don't do `cm.Labels = map[string]string{...}`: that deletes other labels. Create the map if nil, then set your key.

## What you learn
- The single most common controller pattern.
- "Own only your fields": several controllers and humans touch the same object.
- `None` means no API write happened: cheap, idempotent reconciles.

## In Gardener
`pkg/controllerutils/patch.go` (`GetAndCreateOrMergePatch` and friends): Gardener's variants that use patches instead of updates to avoid conflicts.
