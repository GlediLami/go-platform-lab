// Package pods reads pods with the controller-runtime client, like
// `kubectl get pods | grep -v Running` but in code.
package pods

import (
	"context"
	"fmt"
	"slices"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UnhealthyPods returns the sorted names of pods in namespace that are not healthy.
//
// Healthy means: phase Succeeded, or phase Running with every container ready.
func UnhealthyPods(ctx context.Context, c client.Client, namespace string) ([]string, error) {
	podList := &corev1.PodList{}
	if err := c.List(ctx, podList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("listing pods in namespace %q: %w", namespace, err)
	}

	var unhealthy []string
	for _, pod := range podList.Items {
		if !isHealthy(pod) {
			unhealthy = append(unhealthy, pod.Name)
		}
	}
	slices.Sort(unhealthy)
	return unhealthy, nil
}

func isHealthy(pod corev1.Pod) bool {
	switch pod.Status.Phase {
	case corev1.PodSucceeded:
		return true
	case corev1.PodRunning:
		for _, cs := range pod.Status.ContainerStatuses {
			if !cs.Ready {
				return false
			}
		}
		return true
	default:
		return false
	}
}
