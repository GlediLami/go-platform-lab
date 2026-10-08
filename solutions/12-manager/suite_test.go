package manager

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
	"github.com/GlediLami/go-gardener-lab/pkg/labscheme"
)

// This is an integration test: envtest starts a REAL kube-apiserver and etcd
// (no nodes, no kube-controller-manager), installs the CRD, and runs the
// manager against it. Gardener's test/integration tests work the same way.

var (
	testEnv   *envtest.Environment
	restCfg   *rest.Config
	k8sClient client.Client
	fakeCloud *cloud.Fake
	cancel    context.CancelFunc
)

func TestManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "12 Manager Integration Suite")
}

var _ = BeforeSuite(func() {
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		Skip("KUBEBUILDER_ASSETS is not set: run this via `make test EX=12` (it downloads kube-apiserver and etcd)")
	}
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "config", "crd")},
		ErrorIfCRDPathMissing: true,
	}
	var err error
	restCfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(restCfg, client.Options{Scheme: labscheme.Scheme})
	Expect(err).NotTo(HaveOccurred())

	mgr, err := ctrl.NewManager(restCfg, ctrl.Options{
		Scheme:                 labscheme.Scheme,
		Metrics:                metricsserver.Options{BindAddress: "0"},
		HealthProbeBindAddress: "0",
	})
	Expect(err).NotTo(HaveOccurred())

	fakeCloud = cloud.NewFake()
	Expect(Setup(mgr, fakeCloud)).To(Succeed())

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		defer GinkgoRecover()
		Expect(mgr.Start(ctx)).To(Succeed())
	}()
})

var _ = AfterSuite(func() {
	if cancel != nil {
		cancel()
	}
	if testEnv != nil {
		Expect(testEnv.Stop()).To(Succeed())
	}
})
