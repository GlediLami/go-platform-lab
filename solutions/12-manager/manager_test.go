package manager

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
)

var _ = Describe("All controllers together", func() {
	const timeout, interval = 10 * time.Second, 100 * time.Millisecond

	var (
		ctx     context.Context
		ns      *corev1.Namespace
		cluster *v1alpha1.Cluster
	)

	BeforeEach(func() {
		ctx = context.Background()
		ns = &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "team-"}}
		Expect(k8sClient.Create(ctx, ns)).To(Succeed())

		cluster = &v1alpha1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Namespace: ns.Name, Name: "dev"},
			Spec: v1alpha1.ClusterSpec{
				Version: "1.33.2",
				Workers: []v1alpha1.WorkerPool{
					{Name: "pool-a", Minimum: 2, Maximum: 5},
					{Name: "pool-b", Minimum: 1, Maximum: 3},
				},
			},
		}
		Expect(k8sClient.Create(ctx, cluster)).To(Succeed())
	})

	It("reconciles a Cluster end to end", func() {
		By("marking the spec valid")
		Eventually(func(g Gomega) {
			got := &v1alpha1.Cluster{}
			g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cluster), got)).To(Succeed())
			cond := meta.FindStatusCondition(got.Status.Conditions, v1alpha1.ConditionSpecValid)
			g.Expect(cond).NotTo(BeNil())
			g.Expect(cond.Status).To(Equal(metav1.ConditionTrue))
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())

		By("creating the workers Deployment with 3 replicas")
		dep := &appsv1.Deployment{}
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(ctx, client.ObjectKey{Namespace: ns.Name, Name: "dev-workers"}, dep)).To(Succeed())
			g.Expect(ptr.Deref(dep.Spec.Replicas, -1)).To(BeEquivalentTo(3))
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())

		By("creating the network and adding the finalizer")
		netName := "cluster--" + ns.Name + "--dev"
		Eventually(fakeCloud.Networks).WithTimeout(timeout).WithPolling(interval).Should(HaveKey(netName))
		Eventually(func(g Gomega) {
			got := &v1alpha1.Cluster{}
			g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cluster), got)).To(Succeed())
			g.Expect(got.Finalizers).To(ContainElement("lab.example.com/network"))
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())

		By("putting the Deployment back when someone scales it by hand (Owns)")
		patch := client.MergeFrom(dep.DeepCopy())
		dep.Spec.Replicas = ptr.To[int32](10)
		Expect(k8sClient.Patch(ctx, dep, patch)).To(Succeed())
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(dep), dep)).To(Succeed())
			g.Expect(ptr.Deref(dep.Spec.Replicas, -1)).To(BeEquivalentTo(3))
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())

		By("hibernating: workers go to zero")
		Eventually(func(g Gomega) {
			got := &v1alpha1.Cluster{}
			g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(cluster), got)).To(Succeed())
			got.Spec.Hibernated = true
			g.Expect(k8sClient.Update(ctx, got)).To(Succeed())
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())
		Eventually(func(g Gomega) {
			g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(dep), dep)).To(Succeed())
			g.Expect(ptr.Deref(dep.Spec.Replicas, -1)).To(BeEquivalentTo(0))
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())

		By("deleting: network removed, then the Cluster is gone")
		Expect(k8sClient.Delete(ctx, cluster)).To(Succeed())
		Eventually(fakeCloud.Networks).WithTimeout(timeout).WithPolling(interval).ShouldNot(HaveKey(netName))
		Eventually(func() bool {
			err := k8sClient.Get(ctx, client.ObjectKeyFromObject(cluster), &v1alpha1.Cluster{})
			return apierrors.IsNotFound(err)
		}).WithTimeout(timeout).WithPolling(interval).Should(BeTrue())
	})

	It("reports an invalid spec on the status", func() {
		bad := &v1alpha1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Namespace: ns.Name, Name: "bad"},
			Spec: v1alpha1.ClusterSpec{
				Version: "1.33.2",
				Workers: []v1alpha1.WorkerPool{{Name: "ML_Pool", Minimum: 3, Maximum: 1}},
			},
		}
		Expect(k8sClient.Create(ctx, bad)).To(Succeed())

		Eventually(func(g Gomega) {
			got := &v1alpha1.Cluster{}
			g.Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(bad), got)).To(Succeed())
			cond := meta.FindStatusCondition(got.Status.Conditions, v1alpha1.ConditionSpecValid)
			g.Expect(cond).NotTo(BeNil())
			g.Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			g.Expect(cond.Message).To(ContainSubstring("spec.workers[0].maximum"))
		}).WithTimeout(timeout).WithPolling(interval).Should(Succeed())
	})
})
