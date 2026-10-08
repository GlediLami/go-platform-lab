// Package replicas decides how many worker machines each pool should run.
package replicas

import "github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"

// DesiredReplicas returns how many machines a pool should have.
//
// TODO(you):
//   - hibernated -> 0
//   - otherwise keep current inside [pool.Minimum, pool.Maximum]
func DesiredReplicas(pool v1alpha1.WorkerPool, current int32, hibernated bool) int32 {
	return -1
}

// TotalDesired sums DesiredReplicas over all pools. current maps pool name to
// the machines running now; a pool missing from the map has 0 machines.
//
// TODO(you): loop over spec.Workers and reuse DesiredReplicas.
// Hint: reading a missing key from a Go map returns the zero value.
func TotalDesired(spec v1alpha1.ClusterSpec, current map[string]int32) int32 {
	return -1
}
