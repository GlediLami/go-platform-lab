// Package cleanup creates a cloud network for every Cluster and uses a
// finalizer to delete it before the Cluster disappears, like Gardener deletes
// a Shoot's infrastructure before the Shoot object is gone.
package cleanup

import (
	"context"
	"errors"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-gardener-lab/pkg/apis/lab/v1alpha1"
	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
)

const (
	// Finalizer blocks deletion until the network is gone.
	Finalizer = "lab.gardener.cloud/network"
	// NetworkCIDR is the CIDR every Cluster network gets.
	NetworkCIDR = "10.0.0.0/16"
)

// NetworkName is the cloud network name of a Cluster, like the
// shoot--<project>--<name> technical ID in Gardener.
func NetworkName(cluster *v1alpha1.Cluster) string {
	return "cluster--" + cluster.Namespace + "--" + cluster.Name
}

// Reconciler manages the cloud network of each Cluster.
type Reconciler struct {
	Client client.Client
	Cloud  cloud.API
}

// Reconcile ensures the network exists, or deletes it when the Cluster is being deleted.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	cluster := &v1alpha1.Cluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, cluster); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	if !cluster.DeletionTimestamp.IsZero() {
		return r.delete(ctx, cluster)
	}
	return r.reconcile(ctx, cluster)
}

func (r *Reconciler) reconcile(ctx context.Context, cluster *v1alpha1.Cluster) (reconcile.Result, error) {
	// Add the finalizer BEFORE creating anything in the cloud. Otherwise a
	// delete in between would leave an orphaned network behind.
	if !controllerutil.ContainsFinalizer(cluster, Finalizer) {
		patch := client.MergeFrom(cluster.DeepCopy())
		controllerutil.AddFinalizer(cluster, Finalizer)
		if err := r.Client.Patch(ctx, cluster, patch); err != nil {
			return reconcile.Result{}, fmt.Errorf("adding finalizer: %w", err)
		}
	}

	name := NetworkName(cluster)
	if _, err := r.Cloud.GetNetwork(ctx, name); err != nil {
		if !errors.Is(err, cloud.ErrNotFound) {
			return reconcile.Result{}, fmt.Errorf("getting network %q: %w", name, err)
		}
		if _, err := r.Cloud.CreateNetwork(ctx, name, NetworkCIDR); err != nil {
			return reconcile.Result{}, fmt.Errorf("creating network %q: %w", name, err)
		}
		logf.FromContext(ctx).Info("Created network", "network", name)
	}
	return reconcile.Result{}, nil
}

func (r *Reconciler) delete(ctx context.Context, cluster *v1alpha1.Cluster) (reconcile.Result, error) {
	if !controllerutil.ContainsFinalizer(cluster, Finalizer) {
		return reconcile.Result{}, nil // not ours (anymore)
	}

	name := NetworkName(cluster)
	if err := r.Cloud.DeleteNetwork(ctx, name); err != nil && !errors.Is(err, cloud.ErrNotFound) {
		// Keep the finalizer: the Cluster stays until the cloud is clean. The
		// returned error makes controller-runtime retry with backoff.
		return reconcile.Result{}, fmt.Errorf("deleting network %q: %w", name, err)
	}

	patch := client.MergeFrom(cluster.DeepCopy())
	controllerutil.RemoveFinalizer(cluster, Finalizer)
	if err := r.Client.Patch(ctx, cluster, patch); err != nil {
		return reconcile.Result{}, fmt.Errorf("removing finalizer: %w", err)
	}
	logf.FromContext(ctx).Info("Deleted network, released Cluster", "network", name)
	return reconcile.Result{}, nil
}

// SetupWithManager registers the reconciler with a manager.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("cluster-network").
		For(&v1alpha1.Cluster{}).
		Complete(r)
}
