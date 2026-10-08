// Package cleanup creates a cloud network for every Cluster and uses a
// finalizer to delete it before the Cluster disappears, like a managed
// Kubernetes service deletes a cluster's infrastructure before the object is gone.
package cleanup

import (
	"context"
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
	"github.com/GlediLami/go-platform-lab/pkg/cloud"
)

const (
	// Finalizer blocks deletion until the network is gone.
	Finalizer = "lab.example.com/network"
	// NetworkCIDR is the CIDR every Cluster network gets.
	NetworkCIDR = "10.0.0.0/16"
)

// NetworkName is the cloud network name of a Cluster: unique per namespace and name.
//
// TODO(you): "cluster--<namespace>--<name>"
func NetworkName(cluster *v1alpha1.Cluster) string {
	return ""
}

// Reconciler manages the cloud network of each Cluster.
type Reconciler struct {
	Client client.Client
	Cloud  cloud.API
}

// Reconcile ensures the network exists, or deletes it when the Cluster is being deleted.
//
// TODO(you):
//  1. Get the Cluster (not found -> no error).
//  2. Being deleted? (!cluster.DeletionTimestamp.IsZero())
//     - finalizer not there -> nothing to do
//     - r.Cloud.DeleteNetwork; a real error -> return it and KEEP the finalizer;
//     cloud.ErrNotFound counts as success
//     - remove the finalizer (controllerutil.RemoveFinalizer) and Patch the Cluster
//  3. Not being deleted:
//     - add the finalizer FIRST (controllerutil.AddFinalizer + Patch) so a delete
//     can never leave an orphaned network behind
//     - GetNetwork; ErrNotFound -> CreateNetwork(name, NetworkCIDR); other error -> return it
//
// Patch pattern: patch := client.MergeFrom(cluster.DeepCopy()); change cluster; r.Client.Patch(ctx, cluster, patch)
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, errors.New("not implemented")
}

// SetupWithManager registers the reconciler with a manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("cluster-network").
		For(&v1alpha1.Cluster{}).
		Complete(r)
}
