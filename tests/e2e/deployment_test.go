package e2e_test

import (
	"github.com/Emon46/image-cloner/tests/e2e/framework"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
)

var _ = Describe("Deployment", func() {
	var (
		err error
		f   *framework.Invocation
	)
	BeforeEach(func() {
		f = root.Invoke()
	})

	Describe("ImageCloner", func() {
		Context("DeploymentController", func() {
			var (
				deploymentName string
				deployment     *appsv1.Deployment
			)

			BeforeEach(func() {
				deploymentName = f.GetRandomName("")
				deployment = f.GetDeployment(deploymentName)
			})

			AfterEach(func() {
				By("Deleting Deployment")
				err = f.DeleteDeployment(deployment.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Deleting Deployment")
				f.EventuallyDeploymentDeleted(deployment.ObjectMeta).Should(BeTrue())

			})

			It("should create, image clone and delete Deployment successfully", func() {

				By("Creating Deployment")
				err = f.CreateDeployment(deployment)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Image Get cloned")
				f.EventuallyImageClonedForDeploymentToBackupRegistry(deployment.ObjectMeta, registry+"/").Should(BeTrue())
			})
		})
	})
})
