// Package clusterstatus reports validation results on the Cluster's status as
// a condition, the way Gardener reports health on Shoot status.
package clusterstatus

import (
	"context"
	"errors"

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-gardener-lab/pkg/apis/lab/v1alpha1"
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
//
// TODO(you):
//  1. Get the Cluster (not found -> no error).
//  2. Build a metav1.Condition{Type: v1alpha1.ConditionSpecValid, ObservedGeneration: cluster.Generation}
//     - r.Validate(cluster) has errors -> Status False, Reason ReasonInvalidSpec,
//     Message errs.ToAggregate().Error()
//     - no errors -> Status True, Reason ReasonValid, any Message
//  3. meta.SetStatusCondition(&cluster.Status.Conditions, condition)
//     (k8s.io/apimachinery/pkg/api/meta; it replaces a condition of the same Type)
//  4. cluster.Status.ObservedGeneration = cluster.Generation
//  5. Write it with r.Client.Status().Update(ctx, cluster) (NOT r.Client.Update:
//     status is a subresource and a normal update ignores it)
//  6. An invalid spec is not an error to return.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, errors.New("not implemented")
}

// SetupWithManager registers the reconciler with a manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("cluster-status").
		For(&v1alpha1.Cluster{}).
		Complete(r)
}
