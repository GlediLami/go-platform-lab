package cleanup

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-gardener-lab/pkg/apis/lab/v1alpha1"
	"github.com/GlediLami/go-gardener-lab/pkg/cloud"
	"github.com/GlediLami/go-gardener-lab/pkg/labscheme"
)

var _ = Describe("Reconciler", func() {
	var (
		ctx       context.Context
		c         client.Client
		fakeCloud *cloud.Fake
		r         *Reconciler
		cluster   *v1alpha1.Cluster
		request   reconcile.Request
		netName   = "cluster--garden-dev--dev"
	)

	BeforeEach(func() {
		ctx = context.Background()
		cluster = &v1alpha1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Namespace: "garden-dev", Name: "dev"},
			Spec:       v1alpha1.ClusterSpec{Version: "1.33.2"},
		}
		request = reconcile.Request{NamespacedName: client.ObjectKeyFromObject(cluster)}
		fakeCloud = cloud.NewFake()
	})

	JustBeforeEach(func() {
		c = fake.NewClientBuilder().WithScheme(labscheme.Scheme).WithObjects(cluster).Build()
		r = &Reconciler{Client: c, Cloud: fakeCloud}
	})

	doReconcile := func() error {
		_, err := r.Reconcile(ctx, request)
		return err
	}

	getCluster := func() (*v1alpha1.Cluster, error) {
		got := &v1alpha1.Cluster{}
		return got, c.Get(ctx, request.NamespacedName, got)
	}

	It("uses the expected network name", func() {
		Expect(NetworkName(cluster)).To(Equal(netName))
	})

	It("adds the finalizer and creates the network", func() {
		Expect(doReconcile()).To(Succeed())

		got, err := getCluster()
		Expect(err).NotTo(HaveOccurred())
		Expect(got.Finalizers).To(ContainElement(Finalizer))
		Expect(fakeCloud.Networks()).To(HaveKeyWithValue(netName, cloud.Network{Name: netName, CIDR: NetworkCIDR}))
	})

	It("does not create the network twice", func() {
		Expect(doReconcile()).To(Succeed())
		Expect(doReconcile()).To(Succeed())

		creates := 0
		for _, call := range fakeCloud.Calls() {
			if call == "create:"+netName {
				creates++
			}
		}
		Expect(creates).To(Equal(1))
	})

	Context("when the Cluster is deleted", func() {
		JustBeforeEach(func() {
			Expect(doReconcile()).To(Succeed())          // finalizer + network
			Expect(c.Delete(ctx, cluster)).To(Succeed()) // only sets deletionTimestamp: the finalizer blocks removal

			got, err := getCluster()
			Expect(err).NotTo(HaveOccurred(), "the finalizer must keep the Cluster around until cleanup is done")
			Expect(got.DeletionTimestamp.IsZero()).To(BeFalse())
		})

		It("deletes the network, then removes the finalizer so the Cluster goes away", func() {
			Expect(doReconcile()).To(Succeed())

			Expect(fakeCloud.Networks()).NotTo(HaveKey(netName))
			_, err := getCluster()
			Expect(apierrors.IsNotFound(err)).To(BeTrue(), "expected the Cluster to be gone, got err=%v", err)
		})

		It("keeps the finalizer and returns an error when the cloud delete fails", func() {
			fakeCloud.DeleteErr = errors.New("cloud API unavailable")

			Expect(doReconcile()).To(MatchError(ContainSubstring("cloud API unavailable")))

			got, err := getCluster()
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Finalizers).To(ContainElement(Finalizer))
		})

		It("treats an already deleted network as done", func() {
			Expect(fakeCloud.DeleteNetwork(ctx, netName)).To(Succeed()) // someone cleaned up by hand

			Expect(doReconcile()).To(Succeed())
			_, err := getCluster()
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})
	})

	It("ignores Clusters that no longer exist", func() {
		Expect(c.Delete(ctx, cluster)).To(Succeed())
		Expect(doReconcile()).To(Succeed())
	})
})
