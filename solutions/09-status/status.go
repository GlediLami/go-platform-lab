// Package clusterstatus reports validation results on the Cluster's status as
// a condition, the way Kubernetes objects report their state.
package clusterstatus

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
)

// Condition reasons.
const (
	ReasonValid       = "Valid"
	ReasonInvalidSpec = "InvalidSpec"
)

// Reconciler sets the SpecValid condition on Clusters.
type Reconciler struct {
	Client client.Client
	// Validate is injected so the reconciler can be tested with any validation
	// result (dependency injection). In production it is ValidateCluster from exercise 01.
	Validate func(*v1alpha1.Cluster) field.ErrorList
}

// Reconcile validates the Cluster and writes the result to status.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	cluster := &v1alpha1.Cluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, cluster); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	condition := metav1.Condition{
		Type:               v1alpha1.ConditionSpecValid,
		ObservedGeneration: cluster.Generation,
	}
	if errs := r.Validate(cluster); len(errs) > 0 {
		condition.Status = metav1.ConditionFalse
		condition.Reason = ReasonInvalidSpec
		condition.Message = errs.ToAggregate().Error()
	} else {
		condition.Status = metav1.ConditionTrue
		condition.Reason = ReasonValid
		condition.Message = "The spec is valid."
	}

	// SetStatusCondition only bumps LastTransitionTime when Status changes.
	meta.SetStatusCondition(&cluster.Status.Conditions, condition)
	cluster.Status.ObservedGeneration = cluster.Generation

	// Status is a subresource: a normal Update ignores it, Status().Update writes it.
	if err := r.Client.Status().Update(ctx, cluster); err != nil {
		return reconcile.Result{}, fmt.Errorf("updating status of Cluster %s: %w", req.NamespacedName, err)
	}

	// An invalid spec is the user's problem, not a reason to retry.
	return reconcile.Result{}, nil
}

// SetupWithManager registers the reconciler with a manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("cluster-status").
		For(&v1alpha1.Cluster{}).
		Complete(r)
}
