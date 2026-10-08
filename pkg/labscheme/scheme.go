// Package labscheme provides a runtime.Scheme that knows the Kubernetes core
// types and the lab types. A scheme maps Go types to API kinds and back.
package labscheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"github.com/GlediLami/go-gardener-lab/pkg/apis/lab/v1alpha1"
)

// Scheme contains client-go's built-in types plus lab.gardener.cloud/v1alpha1.
var Scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(v1alpha1.AddToScheme(Scheme))
}
