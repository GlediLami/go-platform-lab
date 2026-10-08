// Package minreplicas is a first real reconciler: it makes sure annotated
// Deployments never run fewer replicas than the annotation says.
package minreplicas

import (
	"context"
	"fmt"
	"strconv"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AnnotationMinReplicas holds the minimum replica count, e.g. "3".
const AnnotationMinReplicas = "lab.example.com/min-replicas"

// Reconciler enforces AnnotationMinReplicas on Deployments.
type Reconciler struct {
	Client client.Client
}

// Reconcile is called with the namespace/name of a Deployment that changed.
// It must work no matter how often it is called (level-based, not edge-based).
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	log := logf.FromContext(ctx)

	deployment := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deployment); err != nil {
		// Deleted in the meantime: nothing to do, and no reason to retry.
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	value, ok := deployment.Annotations[AnnotationMinReplicas]
	if !ok {
		return reconcile.Result{}, nil
	}

	minReplicas, err := strconv.ParseInt(value, 10, 32)
	if err != nil || minReplicas < 0 {
		// A user error: retrying won't fix it, so log and don't return an error.
		log.Info("Ignoring invalid annotation", "annotation", AnnotationMinReplicas, "value", value)
		return reconcile.Result{}, nil
	}

	current := ptr.Deref(deployment.Spec.Replicas, 1) // Kubernetes defaults nil to 1
	if current >= int32(minReplicas) {
		return reconcile.Result{}, nil
	}

	patch := client.MergeFrom(deployment.DeepCopy())
	deployment.Spec.Replicas = ptr.To(int32(minReplicas))
	if err := r.Client.Patch(ctx, deployment, patch); err != nil {
		return reconcile.Result{}, fmt.Errorf("scaling Deployment %s up to %d: %w", req.NamespacedName, minReplicas, err)
	}

	log.Info("Scaled up Deployment", "from", current, "to", minReplicas)
	return reconcile.Result{}, nil
}

// SetupWithManager registers the reconciler with a manager (used in exercise 12).
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("min-replicas").
		For(&appsv1.Deployment{}).
		Complete(r)
}
