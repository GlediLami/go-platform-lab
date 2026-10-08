// Package workers makes a Cluster own a Deployment that stands in for its
// worker machines, like the Worker extension creates a MachineDeployment.
package workers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
)

const (
	// LabelCluster points from owned objects back to their Cluster.
	LabelCluster = "lab.example.com/cluster"
	// MachineImage is a tiny image that just sleeps: each pod is a pretend machine.
	MachineImage = "registry.k8s.io/pause:3.10"
)

// DeploymentName returns the name of the workers Deployment of a Cluster.
func DeploymentName(cluster *v1alpha1.Cluster) string {
	return cluster.Name + "-workers"
}

// DesiredReplicas is the sum of all pool minimums, or 0 when hibernated.
func DesiredReplicas(cluster *v1alpha1.Cluster) int32 {
	if cluster.Spec.Hibernated {
		return 0
	}
	var total int32
	for _, pool := range cluster.Spec.Workers {
		total += pool.Minimum
	}
	return total
}

// Reconciler keeps the workers Deployment of each Cluster in sync.
type Reconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile creates or updates the workers Deployment and reports replicas in status.
func (r *Reconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	cluster := &v1alpha1.Cluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, cluster); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	replicas := DesiredReplicas(cluster)
	labels := map[string]string{"app": DeploymentName(cluster), LabelCluster: cluster.Name}

	deployment := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Namespace: cluster.Namespace, Name: DeploymentName(cluster)}}
	if _, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		if deployment.Labels == nil {
			deployment.Labels = map[string]string{}
		}
		for k, v := range labels {
			deployment.Labels[k] = v
		}
		// The selector is immutable: only set it on create.
		if deployment.Spec.Selector == nil {
			deployment.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
		}
		deployment.Spec.Replicas = ptr.To(replicas)
		deployment.Spec.Template.Labels = labels
		deployment.Spec.Template.Spec.Containers = []corev1.Container{{Name: "machine", Image: MachineImage}}

		// Owner reference: Kubernetes garbage-collects the Deployment when the
		// Cluster is deleted, and Owns() in SetupWithManager re-triggers us when
		// someone changes the Deployment.
		return controllerutil.SetControllerReference(cluster, deployment, r.Scheme)
	}); err != nil {
		return reconcile.Result{}, fmt.Errorf("ensuring workers Deployment for Cluster %s: %w", req.NamespacedName, err)
	}

	if cluster.Status.Replicas != replicas {
		cluster.Status.Replicas = replicas
		if err := r.Client.Status().Update(ctx, cluster); err != nil {
			return reconcile.Result{}, fmt.Errorf("updating status of Cluster %s: %w", req.NamespacedName, err)
		}
	}

	return reconcile.Result{}, nil
}

// SetupWithManager registers the reconciler: reconcile on Cluster changes AND
// on changes to Deployments the Cluster owns.
func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		Named("cluster-workers").
		For(&v1alpha1.Cluster{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
