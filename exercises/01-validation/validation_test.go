package validation

import (
	"slices"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/GlediLami/go-gardener-lab/pkg/apis/lab/v1alpha1"
)

// validCluster returns a Cluster that passes validation. Tests change one thing at a time.
func validCluster() *v1alpha1.Cluster {
	return &v1alpha1.Cluster{
		ObjectMeta: metav1.ObjectMeta{Name: "dev"},
		Spec: v1alpha1.ClusterSpec{
			Version: "1.33.2",
			Workers: []v1alpha1.WorkerPool{
				{Name: "pool-a", Minimum: 1, Maximum: 3},
			},
		},
	}
}

// errorKey is "<type> <field>", e.g. "FieldValueRequired spec.version".
func errorKeys(errs field.ErrorList) []string {
	keys := make([]string, 0, len(errs))
	for _, e := range errs {
		keys = append(keys, string(e.Type)+" "+e.Field)
	}
	slices.Sort(keys)
	return keys
}

func TestValidateCluster(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(c *v1alpha1.Cluster)
		want   []string
	}{
		{
			name:   "valid cluster",
			mutate: func(*v1alpha1.Cluster) {},
			want:   []string{},
		},
		{
			name:   "invalid cluster name",
			mutate: func(c *v1alpha1.Cluster) { c.Name = "Dev_Cluster" },
			want:   []string{"FieldValueInvalid metadata.name"},
		},
		{
			name:   "missing version",
			mutate: func(c *v1alpha1.Cluster) { c.Spec.Version = "" },
			want:   []string{"FieldValueRequired spec.version"},
		},
		{
			name:   "no worker pools",
			mutate: func(c *v1alpha1.Cluster) { c.Spec.Workers = nil },
			want:   []string{"FieldValueRequired spec.workers"},
		},
		{
			name:   "invalid pool name",
			mutate: func(c *v1alpha1.Cluster) { c.Spec.Workers[0].Name = "ML_Pool" },
			want:   []string{"FieldValueInvalid spec.workers[0].name"},
		},
		{
			name:   "pool name too long",
			mutate: func(c *v1alpha1.Cluster) { c.Spec.Workers[0].Name = "a-very-long-pool-name" },
			want:   []string{"FieldValueInvalid spec.workers[0].name"},
		},
		{
			name:   "negative minimum",
			mutate: func(c *v1alpha1.Cluster) { c.Spec.Workers[0].Minimum = -1 },
			want:   []string{"FieldValueInvalid spec.workers[0].minimum"},
		},
		{
			name: "maximum below minimum",
			mutate: func(c *v1alpha1.Cluster) {
				c.Spec.Workers[0].Minimum = 3
				c.Spec.Workers[0].Maximum = 1
			},
			want: []string{"FieldValueInvalid spec.workers[0].maximum"},
		},
		{
			name: "duplicate pool names",
			mutate: func(c *v1alpha1.Cluster) {
				c.Spec.Workers = append(c.Spec.Workers, v1alpha1.WorkerPool{Name: "pool-a", Minimum: 1, Maximum: 1})
			},
			want: []string{"FieldValueDuplicate spec.workers[1].name"},
		},
		{
			name: "several problems are all reported",
			mutate: func(c *v1alpha1.Cluster) {
				c.Spec.Version = ""
				c.Spec.Workers[0].Minimum = 3
				c.Spec.Workers[0].Maximum = 1
			},
			want: []string{
				"FieldValueInvalid spec.workers[0].maximum",
				"FieldValueRequired spec.version",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validCluster()
			tt.mutate(c)

			got := errorKeys(ValidateCluster(c))
			if !slices.Equal(got, tt.want) {
				t.Errorf("ValidateCluster() errors\n got: %v\nwant: %v", got, tt.want)
			}
		})
	}
}
