package pods

import (
	"context"
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func pod(namespace, name string, phase corev1.PodPhase, ready ...bool) *corev1.Pod {
	p := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name},
		Status:     corev1.PodStatus{Phase: phase},
	}
	for _, r := range ready {
		p.Status.ContainerStatuses = append(p.Status.ContainerStatuses, corev1.ContainerStatus{Ready: r})
	}
	return p
}

func TestUnhealthyPods(t *testing.T) {
	ns := "shoot--local--local"
	objects := []client.Object{
		pod(ns, "kube-apiserver", corev1.PodRunning, true, true),
		pod(ns, "etcd-main-0", corev1.PodRunning, true, false), // one container not ready
		pod(ns, "etcd-events-0", corev1.PodPending),
		pod(ns, "backup-job", corev1.PodSucceeded),
		pod(ns, "broken", corev1.PodFailed),
		pod("other-namespace", "also-broken", corev1.PodFailed), // must be ignored
	}

	c := fake.NewClientBuilder().WithObjects(objects...).Build()

	got, err := UnhealthyPods(context.Background(), c, ns)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"broken", "etcd-events-0", "etcd-main-0"}
	if !slices.Equal(got, want) {
		t.Errorf("UnhealthyPods() = %v, want %v", got, want)
	}
}

func TestUnhealthyPods_EmptyNamespace(t *testing.T) {
	c := fake.NewClientBuilder().Build()
	got, err := UnhealthyPods(context.Background(), c, "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got %v, want no pods", got)
	}
}
