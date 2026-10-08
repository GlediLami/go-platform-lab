// Package replicas decides how many worker machines each pool should run.
package replicas

import "github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"

// DesiredReplicas returns how many machines a pool should have.
//
// A hibernated cluster has zero machines. Otherwise the current number is kept
// inside [Minimum, Maximum], like the cluster-autoscaler bounds of a node pool.
func DesiredReplicas(pool v1alpha1.WorkerPool, current int32, hibernated bool) int32 {
	if hibernated {
		return 0
	}
	if current < pool.Minimum {
		return pool.Minimum
	}
	if current > pool.Maximum {
		return pool.Maximum
	}
	return current
}

// TotalDesired sums DesiredReplicas over all pools. current maps pool name to
// the machines running now; a pool missing from the map has 0 machines.
func TotalDesired(spec v1alpha1.ClusterSpec, current map[string]int32) int32 {
	var total int32
	for _, pool := range spec.Workers {
		total += DesiredReplicas(pool, current[pool.Name], spec.Hibernated)
	}
	return total
}
