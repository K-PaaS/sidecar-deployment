package v1_test

import (
	"code.cloudfoundry.org/korifi/controllers/controllers/workloads/testutils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("StatefulSet Runner Pod Mutating Webhook", func() {
	var (
		namespace string
		stsPod    *corev1.Pod
	)

	BeforeEach(func() {
		namespace = testutils.PrefixedGUID("ns")
		err := adminClient.Create(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		})
		Expect(err).NotTo(HaveOccurred())

		stsPod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      testutils.PrefixedGUID("pod") + "-1",
				Namespace: namespace,
			},
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{{
					Name:    "init-1",
					Image:   "alpine",
					Command: []string{"sleep", "1234"},
				}},
				Containers: []corev1.Container{{
					Name:    "application",
					Image:   "alpine",
					Command: []string{"sleep", "9876"},
				}},
			},
		}
	})

	JustBeforeEach(func() {
		Expect(adminClient.Create(ctx, stsPod)).To(Succeed())
		lookupKey := client.ObjectKeyFromObject(stsPod)
		Eventually(func(g Gomega) {
			g.Expect(adminClient.Get(ctx, lookupKey, stsPod)).To(Succeed())
		}).Should(Succeed())
	})

	When("the pod has the `korifi.cloudfoundry.org/add-stsr-index: \"true\"` label", func() {
		BeforeEach(func() {
			stsPod.Labels = map[string]string{
				"korifi.cloudfoundry.org/add-stsr-index": "true",
			}
		})

		It("the application container has a CF_INSTANCE_INDEX ENVVAR", func() {
			Expect(stsPod.Labels).To(HaveKeyWithValue("korifi.cloudfoundry.org/add-stsr-index", "true"))
			Expect(stsPod.Spec.Containers[0].Env).NotTo(BeEmpty())
			Expect(stsPod.Spec.Containers[0].Env[0].Name).To(Equal("CF_INSTANCE_INDEX"))
			Expect(stsPod.Spec.Containers[0].Env[0].Value).To(Equal("1"))
		})
	})

	When("the pod does not have the `korifi.cloudfoundry.org/add-stsr-index: \"true\"` label", func() {
		It("the application container has a CF_INSTANCE_INDEX ENVVAR", func() {
			Expect(stsPod.Spec.Containers[0].Env).To(BeEmpty())
		})
	})
})
