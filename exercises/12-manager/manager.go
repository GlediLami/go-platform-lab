// Package manager wires all lab controllers into one controller-runtime
// Manager, the way a real operator binary wires all its controllers.
package manager

import (
	"errors"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/GlediLami/go-platform-lab/pkg/cloud"
)

// Setup registers every lab controller with mgr. cloudAPI is the cloud the
// network controller talks to (a fake in this lab).
//
// TODO(you): register YOUR reconcilers from the earlier exercises, each with
// its SetupWithManager(mgr), and wrap errors with which controller failed:
//
//   - minreplicas.Reconciler   (exercises/08-reconciler)  Client
//   - clusterstatus.Reconciler (exercises/09-status)      Client, Validate: validation.ValidateCluster (exercises/01-validation)
//   - workers.Reconciler       (exercises/10-owned)       Client, Scheme: mgr.GetScheme()
//   - cleanup.Reconciler       (exercises/11-finalizer)   Client, Cloud: cloudAPI
//
// Import them like: workers "github.com/GlediLami/go-platform-lab/exercises/10-owned"
// The client is mgr.GetClient(): it reads from the manager's cache.
func Setup(mgr ctrl.Manager, cloudAPI cloud.API) error {
	return errors.New("not implemented")
}
