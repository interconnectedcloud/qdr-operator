package e2e

import (
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"

	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[interior] Interconnect interior deployment tests", func() {
	f := framework.NewFramework("basic-interior", nil)

	It("Should be able to cerate a default interior deployment", func() {
		testInteriorDefaults(f)
	})

	It("Should be able to scale up a deployment", func() {
		testInteriorScaleUp(f)
	})

	It("Should be able to scale down a deployment", func() {
		testInteriorScaleDown(f)
	})

})

func testInteriorDefaults(f *framework.Framework) {
	By("Creating a default interior interconnect")
	ei, err := f.CreateInterconnect(f.Namespace, 0, func(ei *v1alpha1.Interconnect) {
		ei.Name = "interior-interconnect"
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating a Deployment with 1 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 1, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan defaults")
	Expect(ei.Name).To(Equal("interior-interconnect"))
	Expect(ei.Spec.DeploymentPlan.Size).To(Equal(int32(1)))
	Expect(ei.Spec.DeploymentPlan.Role).To(Equal(v1alpha1.RouterRoleInterior))
	Expect(ei.Spec.DeploymentPlan.Placement).To(Equal(v1alpha1.PlacementAny))

	By("Creating an Interconnect resource in the namespace")
	dep, err := f.GetDeployment("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())
	Expect(*dep.Spec.Replicas).To(Equal(int32(1)))

	By("Creating a service for the interconnect default listeners")
	svc, err := f.GetService("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the owner reference for the service")
	Expect(svc.OwnerReferences[0].APIVersion).To(Equal(framework.GVR))
	Expect(svc.OwnerReferences[0].Name).To(Equal("interior-interconnect"))
	Expect(*svc.OwnerReferences[0].Controller).To(Equal(true))

}

func testInteriorScaleUp(f *framework.Framework) {
	By("Creating an interior interconnect with size 3")
	ei, err := f.CreateInterconnect(f.Namespace, 3, func(ei *v1alpha1.Interconnect) {
		ei.Name = "interior-interconnect"
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating a Deployment with 3 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 3, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan")
	Expect(ei.Name).To(Equal("interior-interconnect"))
	Expect(ei.Spec.DeploymentPlan.Size).To(Equal(int32(3)))
	Expect(ei.Spec.DeploymentPlan.Role).To(Equal(v1alpha1.RouterRoleInterior))
	Expect(ei.Spec.DeploymentPlan.Placement).To(Equal(v1alpha1.PlacementAny))

	By("Scaling the interior interconnect size up")
	ei.Spec.DeploymentPlan.Size = 4
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment has reached 4 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 4, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())
}

func testInteriorScaleDown(f *framework.Framework) {
	By("Creating an interior interconnect with size 3")
	ei, err := f.CreateInterconnect(f.Namespace, 3, func(ei *v1alpha1.Interconnect) {
		ei.Name = "interior-interconnect"
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating a Deployment with 3 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 3, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan")
	Expect(ei.Name).To(Equal("interior-interconnect"))
	Expect(ei.Spec.DeploymentPlan.Size).To(Equal(int32(3)))
	Expect(ei.Spec.DeploymentPlan.Role).To(Equal(v1alpha1.RouterRoleInterior))
	Expect(ei.Spec.DeploymentPlan.Placement).To(Equal(v1alpha1.PlacementAny))

	By("Scaling the interior interconnect size down")
	ei.Spec.DeploymentPlan.Size = 2
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment has reached 2 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 2, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())
}
