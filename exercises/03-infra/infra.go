// Package infra creates and deletes cloud networks idempotently, like the
// Infrastructure actuator of a Gardener provider extension.
package infra

import (
	"context"
	"errors"

	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
)

// ErrCIDRMismatch means the network exists but with a different CIDR.
var ErrCIDRMismatch = errors.New("network exists with a different CIDR")

// EnsureNetwork makes sure a network with the given name and CIDR exists.
// Calling it many times must have the same effect as calling it once.
//
// TODO(you):
//   - GetNetwork first
//   - found with the same CIDR    -> return it, do NOT create
//   - found with a different CIDR -> return an error wrapping ErrCIDRMismatch
//   - errors.Is(err, cloud.ErrNotFound) -> CreateNetwork
//   - any other error             -> return it wrapped with context:
//     fmt.Errorf("getting network %q: %w", name, err)
func EnsureNetwork(ctx context.Context, api cloud.API, name, cidr string) (*cloud.Network, error) {
	return nil, errors.New("not implemented")
}

// DeleteNetwork deletes the network. A network that is already gone is not an error.
//
// TODO(you): ignore errors that wrap cloud.ErrNotFound, wrap all others.
func DeleteNetwork(ctx context.Context, api cloud.API, name string) error {
	return errors.New("not implemented")
}
