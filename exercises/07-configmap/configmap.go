// Package configmap creates or updates an object idempotently with
// controllerutil.CreateOrUpdate, the most common pattern in controllers.
package configmap

import (
	"context"
	"errors"

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
//
// TODO(you):
//   - cm := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}
//   - controllerutil.CreateOrUpdate(ctx, c, cm, func() error { ... })
//   - inside the func: set ONLY what you own (the label and Data);
//     don't replace the whole Labels map, other labels must survive
//   - return the OperationResult (Created / Updated / None)
func EnsureConfigMap(ctx context.Context, c client.Client, namespace, name string, data map[string]string) (controllerutil.OperationResult, error) {
	return controllerutil.OperationResultNone, errors.New("not implemented")
}
