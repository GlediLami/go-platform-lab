// Package configmap creates or updates an object idempotently with
// controllerutil.CreateOrUpdate, the most common pattern in controllers.
package configmap

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// LabelManagedBy marks objects this lab manages.
const (
	LabelManagedBy = "app.kubernetes.io/managed-by"
	ManagedByValue = "go-gardener-lab"
)

// EnsureConfigMap makes the ConfigMap exist with exactly this data and the
// managed-by label. Labels set by others are kept.
func EnsureConfigMap(ctx context.Context, c client.Client, namespace, name string, data map[string]string) (controllerutil.OperationResult, error) {
	cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}

	result, err := controllerutil.CreateOrUpdate(ctx, c, cm, func() error {
		// This function describes the desired state. It runs on a fresh object
		// (create) or on the current object from the cluster (update), so only
		// set the fields you own and leave the rest alone.
		if cm.Labels == nil {
			cm.Labels = map[string]string{}
		}
		cm.Labels[LabelManagedBy] = ManagedByValue
		cm.Data = data
		return nil
	})
	if err != nil {
		return result, fmt.Errorf("ensuring ConfigMap %s/%s: %w", namespace, name, err)
	}
	return result, nil
}
