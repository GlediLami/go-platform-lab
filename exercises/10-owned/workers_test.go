package workers

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
	"github.com/GlediLami/go-platform-lab/pkg/labscheme"
)

var _ = Describe("Reconciler", func() {
	var (
		ctx     context.Context
		c       client.Client
		r       *Reconciler
		cluster *v1alpha1.Cluster
		request reconcile.Request
		depKey  client.ObjectKey
	)

	BeforeEach(func() {
		ctx = context.Background()
		cluster = &v1alpha1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Namespace: "team-dev", Name: "dev", UID: "cluster-uid-1234"},
			Spec: v1alpha1.ClusterSpec{
				Version: "1.33.2",
				Workers: []v1alpha1.WorkerPool{
					{Name: "pool-a", Minimum: 2, Maximum: 5},
					{Name: "pool-b", Minimum: 1, Maximum: 3},
				},
			},
		}
		request = reconcile.Request{NamespacedName: client.ObjectKeyFromObject(cluster)}
		depKey = client.ObjectKey{Namespace: "team-dev", Name: "dev-workers"}
	})

	JustBeforeEach(func() {
		c = fake.NewClientBuilder().
			WithScheme(labscheme.Scheme).
			WithObjects(cluster).
			WithStatusSubresource(&v1alpha1.Cluster{}).
			Build()
		r = &Reconciler{Client: c, Scheme: labscheme.Scheme}
	})

	doReconcile := func() {
		GinkgoHelper()
		_, err := r.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())
	}

	getDeployment := func() *appsv1.Deployment {
		GinkgoHelper()
		dep := &appsv1.Deployment{}
		Expect(c.Get(ctx, depKey, dep)).To(Succeed())
		return dep
	}

	It("creates the workers Deployment with the sum of pool minimums", func() {
		doReconcile()

		dep := getDeployment()
		Expect(ptr.Deref(dep.Spec.Replicas, -1)).To(BeEquivalentTo(3))
		Expect(dep.Spec.Selector).NotTo(BeNil())
		Expect(dep.Spec.Template.Labels).To(HaveKeyWithValue(LabelCluster, "dev"))
		Expect(dep.Spec.Selector.MatchLabels).To(Equal(dep.Spec.Template.Labels))
		Expect(dep.Spec.Template.Spec.Containers).To(HaveLen(1))
	})

	It("makes the Cluster the controller owner of the Deployment", func() {
		doReconcile()

		dep := getDeployment()
		Expect(dep.OwnerReferences).To(HaveLen(1))
		owner := dep.OwnerReferences[0]
		Expect(owner.Kind).To(Equal("Cluster"))
		Expect(owner.APIVersion).To(Equal("lab.example.com/v1alpha1"))
		Expect(owner.Name).To(Equal("dev"))
		Expect(owner.UID).To(BeEquivalentTo("cluster-uid-1234"))
		Expect(ptr.Deref(owner.Controller, false)).To(BeTrue())
	})

	It("reports the replicas in the Cluster status", func() {
		doReconcile()

		got := &v1alpha1.Cluster{}
		Expect(c.Get(ctx, request.NamespacedName, got)).To(Succeed())
		Expect(got.Status.Replicas).To(BeEquivalentTo(3))
	})

	It("follows spec changes", func() {
		doReconcile()

		updated := &v1alpha1.Cluster{}
		Expect(c.Get(ctx, request.NamespacedName, updated)).To(Succeed())
		updated.Spec.Workers[0].Minimum = 4
		Expect(c.Update(ctx, updated)).To(Succeed())

		doReconcile()
		Expect(ptr.Deref(getDeployment().Spec.Replicas, -1)).To(BeEquivalentTo(5))
	})

	It("is idempotent: a second reconcile changes nothing", func() {
		doReconcile()
		before := getDeployment().ResourceVersion

		doReconcile()
		Expect(getDeployment().ResourceVersion).To(Equal(before))
	})

	Context("when hibernated", func() {
		BeforeEach(func() {
			cluster.Spec.Hibernated = true
		})

		It("scales the workers to zero", func() {
			doReconcile()
			Expect(ptr.Deref(getDeployment().Spec.Replicas, -1)).To(BeEquivalentTo(0))
		})
	})

	It("ignores Clusters that no longer exist", func() {
		Expect(c.Delete(ctx, cluster)).To(Succeed())
		_, err := r.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())
	})
})
