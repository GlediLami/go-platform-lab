// Package infra creates and deletes cloud networks idempotently, like the
// Infrastructure actuator of a Gardener provider extension.
package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
)

// ErrCIDRMismatch means the network exists but with a different CIDR.
var ErrCIDRMismatch = errors.New("network exists with a different CIDR")

// EnsureNetwork makes sure a network with the given name and CIDR exists.
// Calling it many times has the same effect as calling it once (idempotent).
func EnsureNetwork(ctx context.Context, api cloud.API, name, cidr string) (*cloud.Network, error) {
	existing, err := api.GetNetwork(ctx, name)
	switch {
	case err == nil:
		if existing.CIDR != cidr {
			return nil, fmt.Errorf("network %q has CIDR %q, want %q: %w", name, existing.CIDR, cidr, ErrCIDRMismatch)
		}
		return existing, nil
	case errors.Is(err, cloud.ErrNotFound):
		created, err := api.CreateNetwork(ctx, name, cidr)
		if err != nil {
			return nil, fmt.Errorf("creating network %q: %w", name, err)
		}
		return created, nil
	default:
		return nil, fmt.Errorf("getting network %q: %w", name, err)
	}
}

// DeleteNetwork deletes the network. A network that is already gone is not an error.
func DeleteNetwork(ctx context.Context, api cloud.API, name string) error {
	if err := api.DeleteNetwork(ctx, name); err != nil && !errors.Is(err, cloud.ErrNotFound) {
		return fmt.Errorf("deleting network %q: %w", name, err)
	}
	return nil
}
