// Package manager wires all lab controllers into one controller-runtime
// Manager, like cmd/gardenlet wires all gardenlet controllers.
package manager

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
	validation "github.com/GlediLami/go-gardener-lab/solutions/01-validation"
	minreplicas "github.com/GlediLami/go-gardener-lab/solutions/08-reconciler"
	clusterstatus "github.com/GlediLami/go-gardener-lab/solutions/09-status"
	workers "github.com/GlediLami/go-gardener-lab/solutions/10-owned"
	cleanup "github.com/GlediLami/go-gardener-lab/solutions/11-finalizer"
)

// Setup registers every lab controller with mgr. cloudAPI is the cloud the
// network controller talks to (a fake in this lab).
func Setup(mgr ctrl.Manager, cloudAPI cloud.API) error {
	c := mgr.GetClient()

	if err := (&minreplicas.Reconciler{Client: c}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("setting up min-replicas controller: %w", err)
	}
	if err := (&clusterstatus.Reconciler{Client: c, Validate: validation.ValidateCluster}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("setting up cluster-status controller: %w", err)
	}
	if err := (&workers.Reconciler{Client: c, Scheme: mgr.GetScheme()}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("setting up cluster-workers controller: %w", err)
	}
	if err := (&cleanup.Reconciler{Client: c, Cloud: cloudAPI}).SetupWithManager(mgr); err != nil {
		return fmt.Errorf("setting up cluster-network controller: %w", err)
	}
	return nil
}
