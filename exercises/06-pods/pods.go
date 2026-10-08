// Package pods reads pods with the controller-runtime client, like
// `kubectl get pods | grep -v Running` but in code.
package pods

import (
	"context"
	"errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UnhealthyPods returns the sorted names of pods in namespace that are not healthy.
//
// Healthy means: phase Succeeded, or phase Running with every container ready.
//
// TODO(you):
//   - list pods: podList := &corev1.PodList{}; c.List(ctx, podList, client.InNamespace(namespace))
//     (corev1 is k8s.io/api/core/v1)
//   - wrap a List error with context
//   - check pod.Status.Phase and pod.Status.ContainerStatuses[i].Ready
//   - sort the result (slices.Sort)
func UnhealthyPods(ctx context.Context, c client.Client, namespace string) ([]string, error) {
	return nil, errors.New("not implemented")
}
