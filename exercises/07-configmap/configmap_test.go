package configmap

import (
	"context"
	"maps"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const ns, name = "garden", "shoot-info"

func get(t *testing.T, c client.Client) *corev1.ConfigMap {
	t.Helper()
	cm := &corev1.ConfigMap{}
	if err := c.Get(context.Background(), client.ObjectKey{Namespace: ns, Name: name}, cm); err != nil {
		t.Fatalf("getting ConfigMap: %v", err)
	}
	return cm
}

func TestEnsureConfigMap_CreateThenNoop(t *testing.T) {
	ctx := context.Background()
	c := fake.NewClientBuilder().Build()
	data := map[string]string{"version": "1.33.2"}

	res, err := EnsureConfigMap(ctx, c, ns, name, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != controllerutil.OperationResultCreated {
		t.Errorf("first call result = %q, want %q", res, controllerutil.OperationResultCreated)
	}

	cm := get(t, c)
	if !maps.Equal(cm.Data, data) {
		t.Errorf("data = %v, want %v", cm.Data, data)
	}
	if cm.Labels[LabelManagedBy] != ManagedByValue {
		t.Errorf("missing label %s=%s, labels: %v", LabelManagedBy, ManagedByValue, cm.Labels)
	}

	res, err = EnsureConfigMap(ctx, c, ns, name, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != controllerutil.OperationResultNone {
		t.Errorf("second identical call result = %q, want %q (nothing changed)", res, controllerutil.OperationResultNone)
	}
}

func TestEnsureConfigMap_UpdatesDataKeepsForeignLabels(t *testing.T) {
	existing := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Namespace: ns, Name: name, Labels: map[string]string{"team": "ske"}},
		Data:       map[string]string{"version": "1.32.0", "old-key": "x"},
	}
	c := fake.NewClientBuilder().WithObjects(existing).Build()
	data := map[string]string{"version": "1.33.2"}

	res, err := EnsureConfigMap(context.Background(), c, ns, name, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != controllerutil.OperationResultUpdated {
		t.Errorf("result = %q, want %q", res, controllerutil.OperationResultUpdated)
	}

	cm := get(t, c)
	if !maps.Equal(cm.Data, data) {
		t.Errorf("data = %v, want exactly %v", cm.Data, data)
	}
	if cm.Labels["team"] != "ske" {
		t.Errorf("label team=ske was removed; keep labels you don't own. labels: %v", cm.Labels)
	}
	if cm.Labels[LabelManagedBy] != ManagedByValue {
		t.Errorf("missing label %s=%s", LabelManagedBy, ManagedByValue)
	}
}
