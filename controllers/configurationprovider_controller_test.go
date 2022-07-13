package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	appconfigv1alpha1 "appconfig/sync/api/v1alpha1"
)

var _ = Describe("CronJob controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		CronjobName      = "test-cronjob"
		CronjobNamespace = "default"
		JobName          = "test-job"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When updating CronJob Status", func() {
		It("Should increase CronJob Status.Active count when new Jobs are created", func() {
			By("By creating a new CronJob")
			ctx := context.Background()
			cronJob := &appconfigv1alpha1.ConfigurationProvider{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "appconfig.kubernetes.config/v1alpha1",
					Kind:       "ConfigurationProvider",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      CronjobName,
					Namespace: CronjobNamespace,
				},
				Spec: appconfigv1alpha1.ConfigurationProviderSpec{
					Endpoint:      "junbchenconfig",
					ConfigMapName: "tobecreated",
					ClientSecret:  "hide",
					ClientId:      "hide",
					TenantId:      "hide",
				},
			}
			Expect(k8sClient.Create(ctx, cronJob)).Should(Succeed())

			cronjobLookupKey := types.NamespacedName{Name: CronjobName, Namespace: CronjobNamespace}
			createdCronjob := &appconfigv1alpha1.ConfigurationProvider{}

			// We'll need to retry getting this newly created CronJob, given that creation may not immediately happen.
			Eventually(func() bool {
				err := k8sClient.Get(ctx, cronjobLookupKey, createdCronjob)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Let's make sure our Schedule string value was properly converted/handled.
			Expect(createdCronjob.Spec.Endpoint).Should(Equal("junbchenconfig"))

			configMapToCreate := &corev1.ConfigMap{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "tobecreated", Namespace: CronjobNamespace}, configMapToCreate)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())

			Expect(configMapToCreate.Data["message"]).Should(Equal("hello_from_appconfig"))

		})
	})
})
