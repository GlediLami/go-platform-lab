// Package minreplicas is a first real reconciler: it makes sure annotated
// Deployments never run fewer replicas than the annotation says.
package minreplicas

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// AnnotationMinReplicas holds the minimum replica count, e.g. "3".
const AnnotationMinReplicas = "lab.example.com/min-replicas"

// Reconciler enforces AnnotationMinReplicas on Deployments.
type Reconciler struct {
	Client client.Client
}

// Reconcile is called with the namespace/name of a Deployment that changed.
//
// TODO(you):
//  1. Get the Deployment. Not found -> return no error (client.IgnoreNotFound).
//  2. No annotation -> done.
//  3. Parse it (strconv.ParseInt). Invalid or negative -> log, return NO error
//     (a user typo won't fix itself by retrying). Logger: logf.FromContext(ctx)
//     from sigs.k8s.io/controller-runtime/pkg/log
//  4. Current replicas: ptr.Deref(deployment.Spec.Replicas, 1) (k8s.io/utils/ptr)
//  5. If current < minimum: patch := client.MergeFrom(deployment.DeepCopy()),
//     set Spec.Replicas = ptr.To(int32(min)), r.Client.Patch(ctx, deployment, patch)
//  6. Return a wrapped error only for real failures (API errors).
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, errors.New("not implemented")
}

// SetupWithManager registers the reconciler with a manager (used in exercise 12).
// Already done for you: read it.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("min-replicas").
		For(&appsv1.Deployment{}).
		Complete(r)
}
