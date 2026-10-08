// Package v1alpha1 contains a tiny, Shoot-like API used by the lab exercises.
//
// In Gardener, API types live in pkg/apis/core/v1beta1 and their DeepCopy
// functions are generated with controller-gen / deepcopy-gen. Here they are
// written by hand (see zz_deepcopy.go) so you can read what the generator does.
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupName is the API group of the lab types.
const GroupName = "lab.gardener.cloud"

var (
	// SchemeGroupVersion is the group and version used to register these objects.
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	// SchemeBuilder collects the functions that add the lab types to a scheme.
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme adds the lab types to a scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &Cluster{}, &ClusterList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// Condition types set on Cluster status.
const (
	// ConditionSpecValid tells whether the Cluster spec passed validation.
	ConditionSpecValid = "SpecValid"
)

// Cluster is a simplified Shoot: a Kubernetes version, worker pools and hibernation.
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// ClusterSpec is the desired state of a Cluster.
type ClusterSpec struct {
	// Version is the Kubernetes version, e.g. "1.33.2".
	Version string `json:"version"`
	// Workers are the worker pools of the cluster.
	Workers []WorkerPool `json:"workers"`
	// Hibernated scales all workers to zero when true.
	Hibernated bool `json:"hibernated,omitempty"`
}

// WorkerPool is a group of identical worker machines.
type WorkerPool struct {
	// Name of the pool. Must be a DNS label of at most 15 characters.
	Name string `json:"name"`
	// Minimum number of machines.
	Minimum int32 `json:"minimum"`
	// Maximum number of machines.
	Maximum int32 `json:"maximum"`
}

// ClusterStatus is the observed state of a Cluster.
type ClusterStatus struct {
	// ObservedGeneration is the generation the controller last acted on.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Conditions describe the current state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// Replicas is the number of worker machines the controller asked for.
	Replicas int32 `json:"replicas,omitempty"`
}

// ClusterList is a list of Clusters.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Cluster `json:"items"`
}
