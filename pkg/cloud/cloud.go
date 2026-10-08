// Package cloud is a pretend cloud provider API, plus an in-memory fake.
//
// Gardener never talks to a cloud directly from its core. Provider extensions
// (gardener-extension-provider-aws, -openstack, -stackit, ...) wrap the cloud
// SDK behind small interfaces like API below, so that controllers can be unit
// tested against a fake instead of a real cloud.
package cloud

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// ErrNotFound is returned when a resource does not exist.
var ErrNotFound = errors.New("not found")

// Network is a cloud network (think: VPC or OpenStack network).
type Network struct {
	Name string
	CIDR string
}

// API is the part of a cloud provider SDK the lab controllers need.
type API interface {
	// GetNetwork returns the network or an error wrapping ErrNotFound.
	GetNetwork(ctx context.Context, name string) (*Network, error)
	// CreateNetwork creates a network.
	CreateNetwork(ctx context.Context, name, cidr string) (*Network, error)
	// DeleteNetwork deletes a network or returns an error wrapping ErrNotFound.
	DeleteNetwork(ctx context.Context, name string) error
}

// Fake is an in-memory API that records every call. Safe for concurrent use.
type Fake struct {
	mu       sync.Mutex
	networks map[string]Network
	calls    []string

	// Set these to make the next calls fail.
	GetErr    error
	CreateErr error
	DeleteErr error
}

var _ API = &Fake{}

// NewFake returns an empty fake cloud.
func NewFake(existing ...Network) *Fake {
	f := &Fake{networks: map[string]Network{}}
	for _, n := range existing {
		f.networks[n.Name] = n
	}
	return f
}

// GetNetwork implements API.
func (f *Fake) GetNetwork(_ context.Context, name string) (*Network, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "get:"+name)
	if f.GetErr != nil {
		return nil, f.GetErr
	}
	n, ok := f.networks[name]
	if !ok {
		return nil, fmt.Errorf("network %q: %w", name, ErrNotFound)
	}
	return &n, nil
}

// CreateNetwork implements API.
func (f *Fake) CreateNetwork(_ context.Context, name, cidr string) (*Network, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "create:"+name)
	if f.CreateErr != nil {
		return nil, f.CreateErr
	}
	if _, ok := f.networks[name]; ok {
		return nil, fmt.Errorf("network %q already exists", name)
	}
	n := Network{Name: name, CIDR: cidr}
	f.networks[name] = n
	return &n, nil
}

// DeleteNetwork implements API.
func (f *Fake) DeleteNetwork(_ context.Context, name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, "delete:"+name)
	if f.DeleteErr != nil {
		return f.DeleteErr
	}
	if _, ok := f.networks[name]; !ok {
		return fmt.Errorf("network %q: %w", name, ErrNotFound)
	}
	delete(f.networks, name)
	return nil
}

// Calls returns the recorded calls, e.g. ["get:a", "create:a"].
func (f *Fake) Calls() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.calls...)
}

// Networks returns a copy of all networks.
func (f *Fake) Networks() map[string]Network {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make(map[string]Network, len(f.networks))
	for k, v := range f.networks {
		out[k] = v
	}
	return out
}
