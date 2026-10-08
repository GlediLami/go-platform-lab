package replicas

import (
	"testing"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
)

func TestDesiredReplicas(t *testing.T) {
	pool := v1alpha1.WorkerPool{Name: "pool-a", Minimum: 2, Maximum: 5}

	tests := []struct {
		name       string
		current    int32
		hibernated bool
		want       int32
	}{
		{name: "below minimum is raised to minimum", current: 0, want: 2},
		{name: "inside the range stays", current: 3, want: 3},
		{name: "exactly minimum stays", current: 2, want: 2},
		{name: "exactly maximum stays", current: 5, want: 5},
		{name: "above maximum is lowered to maximum", current: 9, want: 5},
		{name: "hibernated is always zero", current: 3, hibernated: true, want: 0},
		// TODO(you): add at least two more cases you think are interesting,
		// e.g. a pool with Minimum 0, or hibernated with current 0.
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DesiredReplicas(pool, tt.current, tt.hibernated); got != tt.want {
				t.Errorf("DesiredReplicas(%+v, %d, %t) = %d, want %d", pool, tt.current, tt.hibernated, got, tt.want)
			}
		})
	}
}

func TestTotalDesired(t *testing.T) {
	spec := v1alpha1.ClusterSpec{
		Workers: []v1alpha1.WorkerPool{
			{Name: "small", Minimum: 1, Maximum: 2},
			{Name: "big", Minimum: 3, Maximum: 10},
		},
	}

	tests := []struct {
		name       string
		current    map[string]int32
		hibernated bool
		want       int32
	}{
		{name: "nothing running yet", current: nil, want: 1 + 3},
		{name: "some running", current: map[string]int32{"small": 2, "big": 5}, want: 2 + 5},
		{name: "too many running", current: map[string]int32{"small": 7, "big": 50}, want: 2 + 10},
		{name: "hibernated", current: map[string]int32{"small": 2, "big": 5}, hibernated: true, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := spec
			s.Hibernated = tt.hibernated
			if got := TotalDesired(s, tt.current); got != tt.want {
				t.Errorf("TotalDesired() = %d, want %d", got, tt.want)
			}
		})
	}
}
