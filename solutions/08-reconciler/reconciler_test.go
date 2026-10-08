package minreplicas

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func deployment(replicas *int32, annotations map[string]string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Namespace: "default", Name: "web", Annotations: annotations},
		Spec:       appsv1.DeploymentSpec{Replicas: replicas},
	}
}

var key = types.NamespacedName{Namespace: "default", Name: "web"}

func reconcileAndGet(t *testing.T, objs ...client.Object) (*appsv1.Deployment, error) {
	t.Helper()
	c := fake.NewClientBuilder().WithObjects(objs...).Build()
	r := &Reconciler{Client: c}

	_, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: key})

	got := &appsv1.Deployment{}
	if getErr := c.Get(context.Background(), key, got); getErr != nil {
		return nil, err
	}
	return got, err
}

func TestReconcile(t *testing.T) {
	tests := []struct {
		name         string
		replicas     *int32
		annotations  map[string]string
		wantReplicas int32
	}{
		{name: "scales up below minimum", replicas: ptr.To[int32](1), annotations: map[string]string{AnnotationMinReplicas: "3"}, wantReplicas: 3},
		{name: "leaves higher count alone", replicas: ptr.To[int32](5), annotations: map[string]string{AnnotationMinReplicas: "3"}, wantReplicas: 5},
		{name: "equal is fine", replicas: ptr.To[int32](3), annotations: map[string]string{AnnotationMinReplicas: "3"}, wantReplicas: 3},
		{name: "nil replicas counts as 1", replicas: nil, annotations: map[string]string{AnnotationMinReplicas: "2"}, wantReplicas: 2},
		{name: "no annotation, no change", replicas: ptr.To[int32](1), annotations: nil, wantReplicas: 1},
		{name: "invalid annotation is ignored", replicas: ptr.To[int32](1), annotations: map[string]string{AnnotationMinReplicas: "lots"}, wantReplicas: 1},
		{name: "negative annotation is ignored", replicas: ptr.To[int32](1), annotations: map[string]string{AnnotationMinReplicas: "-2"}, wantReplicas: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := reconcileAndGet(t, deployment(tt.replicas, tt.annotations))
			if err != nil {
				t.Fatalf("Reconcile returned error: %v (user errors must not be returned as errors)", err)
			}
			if r := ptr.Deref(got.Spec.Replicas, 1); r != tt.wantReplicas {
				t.Errorf("replicas = %d, want %d", r, tt.wantReplicas)
			}
		})
	}
}

func TestReconcile_DeploymentGone(t *testing.T) {
	c := fake.NewClientBuilder().Build()
	r := &Reconciler{Client: c}
	res, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: key})
	if err != nil {
		t.Fatalf("a deleted Deployment must not be an error, got: %v", err)
	}
	if !res.IsZero() {
		t.Errorf("result = %+v, want no requeue", res)
	}
}

func TestReconcile_IsIdempotent(t *testing.T) {
	c := fake.NewClientBuilder().WithObjects(deployment(ptr.To[int32](1), map[string]string{AnnotationMinReplicas: "3"})).Build()
	r := &Reconciler{Client: c}

	for range 3 {
		if _, err := r.Reconcile(context.Background(), reconcile.Request{NamespacedName: key}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	got := &appsv1.Deployment{}
	_ = c.Get(context.Background(), key, got)
	if r := ptr.Deref(got.Spec.Replicas, 1); r != 3 {
		t.Errorf("replicas = %d, want 3", r)
	}
}
