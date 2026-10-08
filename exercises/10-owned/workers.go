// Package workers makes a Cluster own a Deployment that stands in for its
// worker machines, like the Worker extension creates a MachineDeployment.
package workers

import (
	"context"
	"errors"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-gardener-lab/pkg/apis/lab/v1alpha1"
)

const (
	// LabelCluster points from owned objects back to their Cluster.
	LabelCluster = "lab.gardener.cloud/cluster"
	// MachineImage is a tiny image that just sleeps: each pod is a pretend machine.
	MachineImage = "registry.k8s.io/pause:3.10"
)

// DeploymentName returns the name of the workers Deployment of a Cluster.
func DeploymentName(cluster *v1alpha1.Cluster) string {
	return cluster.Name + "-workers"
}

// DesiredReplicas is the sum of all pool minimums, or 0 when hibernated.
//
// TODO(you)
func DesiredReplicas(cluster *v1alpha1.Cluster) int32 {
	return -1
}

// Reconciler keeps the workers Deployment of each Cluster in sync.
type Reconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile creates or updates the workers Deployment and reports replicas in status.
//
// TODO(you):
//  1. Get the Cluster (not found -> no error).
//  2. labels := {"app": DeploymentName(cluster), LabelCluster: cluster.Name}
//  3. controllerutil.CreateOrUpdate a Deployment named DeploymentName(cluster)
//     in the Cluster's namespace. In the mutate func:
//     - add the labels to deployment.Labels (keep others)
//     - Spec.Selector only if nil (it's immutable after create): MatchLabels: labels
//     - Spec.Replicas = ptr.To(DesiredReplicas(cluster))
//     - Spec.Template.Labels = labels
//     - one container {Name: "machine", Image: MachineImage}
//     - return controllerutil.SetControllerReference(cluster, deployment, r.Scheme)
//  4. If cluster.Status.Replicas differs, set it and r.Client.Status().Update.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	return reconcile.Result{}, errors.New("not implemented")
}

// SetupWithManager registers the reconciler: reconcile on Cluster changes AND
// on changes to Deployments the Cluster owns. Read it: Owns() is what makes
// "someone scaled my Deployment" trigger a reconcile.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("cluster-workers").
		For(&v1alpha1.Cluster{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
