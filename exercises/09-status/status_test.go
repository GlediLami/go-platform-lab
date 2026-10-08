package clusterstatus

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/GlediLami/go-platform-lab/pkg/apis/lab/v1alpha1"
	"github.com/GlediLami/go-platform-lab/pkg/labscheme"
)

var _ = Describe("Reconciler", func() {
	var (
		ctx          context.Context
		c            client.Client
		r            *Reconciler
		cluster      *v1alpha1.Cluster
		validateErrs field.ErrorList
		request      reconcile.Request
	)

	BeforeEach(func() {
		ctx = context.Background()
		validateErrs = nil
		cluster = &v1alpha1.Cluster{
			ObjectMeta: metav1.ObjectMeta{Namespace: "team-dev", Name: "dev", Generation: 2},
			Spec: v1alpha1.ClusterSpec{
				Version: "1.33.2",
				Workers: []v1alpha1.WorkerPool{{Name: "pool-a", Minimum: 1, Maximum: 3}},
			},
		}
		request = reconcile.Request{NamespacedName: client.ObjectKeyFromObject(cluster)}
	})

	// JustBeforeEach runs after all BeforeEach blocks, so specs can change the
	// inputs (cluster, validateErrs) first.
	JustBeforeEach(func() {
		c = fake.NewClientBuilder().
			WithScheme(labscheme.Scheme).
			WithObjects(cluster).
			WithStatusSubresource(&v1alpha1.Cluster{}).
			Build()
		r = &Reconciler{
			Client:   c,
			Validate: func(*v1alpha1.Cluster) field.ErrorList { return validateErrs },
		}
	})

	reconcileAndGet := func() *v1alpha1.Cluster {
		GinkgoHelper()
		_, err := r.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())

		got := &v1alpha1.Cluster{}
		Expect(c.Get(ctx, request.NamespacedName, got)).To(Succeed())
		return got
	}

	It("sets SpecValid=True for a valid spec", func() {
		got := reconcileAndGet()

		cond := meta.FindStatusCondition(got.Status.Conditions, v1alpha1.ConditionSpecValid)
		Expect(cond).NotTo(BeNil())
		Expect(cond.Status).To(Equal(metav1.ConditionTrue))
		Expect(cond.Reason).To(Equal(ReasonValid))
		Expect(cond.LastTransitionTime.IsZero()).To(BeFalse())
	})

	It("records the generation it acted on", func() {
		got := reconcileAndGet()

		Expect(got.Status.ObservedGeneration).To(Equal(got.Generation))
		cond := meta.FindStatusCondition(got.Status.Conditions, v1alpha1.ConditionSpecValid)
		Expect(cond).NotTo(BeNil())
		Expect(cond.ObservedGeneration).To(Equal(got.Generation))
	})

	Context("when validation fails", func() {
		BeforeEach(func() {
			validateErrs = field.ErrorList{
				field.Required(field.NewPath("spec", "version"), "must specify a Kubernetes version"),
			}
		})

		It("sets SpecValid=False with the errors as message and no error", func() {
			got := reconcileAndGet()

			cond := meta.FindStatusCondition(got.Status.Conditions, v1alpha1.ConditionSpecValid)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionFalse))
			Expect(cond.Reason).To(Equal(ReasonInvalidSpec))
			Expect(cond.Message).To(ContainSubstring("spec.version"))
		})

		It("flips to True once the spec is fixed", func() {
			reconcileAndGet()

			validateErrs = nil
			got := reconcileAndGet()

			cond := meta.FindStatusCondition(got.Status.Conditions, v1alpha1.ConditionSpecValid)
			Expect(cond).NotTo(BeNil())
			Expect(cond.Status).To(Equal(metav1.ConditionTrue))
			Expect(got.Status.Conditions).To(HaveLen(1), "update the condition, don't append a second one")
		})
	})

	It("ignores Clusters that no longer exist", func() {
		Expect(c.Delete(ctx, cluster)).To(Succeed())

		_, err := r.Reconcile(ctx, request)
		Expect(err).NotTo(HaveOccurred())
	})
})
