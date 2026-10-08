// Package validation validates Cluster objects the way an API server validates resources.
package validation

import (
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
)

// MaxPoolNameLength is the longest allowed worker pool name.
const MaxPoolNameLength = 15

// ValidateCluster validates the whole Cluster object.
//
// TODO(you):
//   - metadata.name must be a DNS-1123 label (use validation.IsDNS1123Label from
//     k8s.io/apimachinery/pkg/util/validation; it returns a list of messages)
//   - then validate the spec with ValidateClusterSpec, using the path "spec"
func ValidateCluster(cluster *v1alpha1.Cluster) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}

// ValidateClusterSpec validates the spec. fldPath is the path of the spec itself.
//
// TODO(you): report ALL problems, not just the first one:
//   - version must not be empty                            -> field.Required
//   - at least one worker pool                             -> field.Required
//   - each pool name is a DNS-1123 label                   -> field.Invalid
//   - each pool name is at most MaxPoolNameLength chars     -> field.Invalid
//   - minimum >= 0                                         -> field.Invalid
//   - maximum >= minimum                                   -> field.Invalid
//   - pool names are unique                                -> field.Duplicate
//
// Paths look like spec.workers[0].name: fldPath.Child("workers").Index(0).Child("name")
func ValidateClusterSpec(spec *v1alpha1.ClusterSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	return allErrs
}
