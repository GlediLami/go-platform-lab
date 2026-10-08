// Package validation validates Cluster objects the way an API server validates resources.
package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
)

// MaxPoolNameLength is the longest allowed worker pool name.
const MaxPoolNameLength = 15

// ValidateCluster validates the whole Cluster object.
func ValidateCluster(cluster *v1alpha1.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}

	namePath := field.NewPath("metadata", "name")
	for _, msg := range validation.IsDNS1123Label(cluster.Name) {
		allErrs = append(allErrs, field.Invalid(namePath, cluster.Name, msg))
	}

	allErrs = append(allErrs, ValidateClusterSpec(&cluster.Spec, field.NewPath("spec"))...)
	return allErrs
}

// ValidateClusterSpec validates the spec. fldPath is the path of the spec itself.
func ValidateClusterSpec(spec *v1alpha1.ClusterSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if spec.Version == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("version"), "must specify a Kubernetes version"))
	}

	workersPath := fldPath.Child("workers")
	if len(spec.Workers) == 0 {
		allErrs = append(allErrs, field.Required(workersPath, "must specify at least one worker pool"))
	}

	seen := map[string]bool{}
	for i, pool := range spec.Workers {
		idxPath := workersPath.Index(i)
		allErrs = append(allErrs, validateWorkerPool(pool, idxPath)...)

		if seen[pool.Name] {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), pool.Name))
		}
		seen[pool.Name] = true
	}

	return allErrs
}

func validateWorkerPool(pool v1alpha1.WorkerPool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	namePath := fldPath.Child("name")
	for _, msg := range validation.IsDNS1123Label(pool.Name) {
		allErrs = append(allErrs, field.Invalid(namePath, pool.Name, msg))
	}
	if len(pool.Name) > MaxPoolNameLength {
		allErrs = append(allErrs, field.Invalid(namePath, pool.Name, fmt.Sprintf("must be at most %d characters", MaxPoolNameLength)))
	}

	if pool.Minimum < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("minimum"), pool.Minimum, "must be greater than or equal to 0"))
	}
	if pool.Maximum < pool.Minimum {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("maximum"), pool.Maximum, "must be greater than or equal to minimum"))
	}

	return allErrs
}
