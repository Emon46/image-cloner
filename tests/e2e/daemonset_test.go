package e2e_test

import (
	"github.com/Emon46/image-cloner/tests/e2e/framework"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

var _ = Describe("DaemonSet", func() {
	var (
		err error
		f   *framework.Invocation
	)
	BeforeEach(func() {
		f = root.Invoke()
	})

	Describe("ImageCloner-DaemonSet", func() {
		Context("DaemonSetController", func() {
			var (
				daemonSetName string
				daemonSet     *appsv1.DaemonSet
			)

			BeforeEach(func() {
				daemonSetName = f.GetRandomName("")
				daemonSet = f.GetDaemonSet(daemonSetName)
			})

			AfterEach(func() {
				By("Deleting DaemonSet")
				err = f.DeleteDaemonSet(daemonSet.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Deleting DaemonSet")
				f.EventuallyDaemonSetDeleted(daemonSet.ObjectMeta).Should(BeTrue())

			})

			It("should create, image clone and delete DaemonSet successfully", func() {
				By("Creating Secret")

				By("Creating DaemonSet")
				err = f.CreateDaemonset(daemonSet)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Image Get cloned")
				f.EventuallyImageClonedForDaemonSetToBackupRegistry(daemonSet.ObjectMeta, registry+"/").Should(BeTrue())
			})
		})
	})
})
