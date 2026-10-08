package clusterstatus

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Many Kubernetes projects test with Ginkgo (Describe/It) and Gomega (Expect(...).To(...)).
// One TestXxx function hands control to Ginkgo, which runs every spec in the package.
func TestClusterStatus(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "09 Cluster status Suite")
}
