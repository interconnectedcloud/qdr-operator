package e2e

import (
	"time"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[edge] Interconnect edge deployment tests", func() {

	f := framework.NewFramework("basic-edge", nil)

	It("Should be able to create a default edge deployment", func() {
		testEdgeDefaults(f)
	})

})

func testEdgeDefaults(f *framework.Framework) {
	By("Creating an edge interconnect with default size")
	ei, err := f.CreateInterconnect(f.Namespace, 0, func(ei *v1alpha1.Interconnect) {
		ei.Name = "edge-interconnect"
		ei.Spec.DeploymentPlan.Role = "edge"
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating a Deployment with 1 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "edge-interconnect", 1, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect("edge-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan")
	Expect(ei.Name).To(Equal("edge-interconnect"))
	Expect(ei.Spec.DeploymentPlan.Size).To(Equal(int32(1)))
	Expect(ei.Spec.DeploymentPlan.Role).To(Equal(v1alpha1.RouterRoleType("edge")))
	Expect(ei.Spec.DeploymentPlan.Placement).To(Equal(v1alpha1.PlacementAny))

	By("Creating an Interconnect resource in the namespace")
	dep, err := f.GetDeployment("edge-interconnect")
	Expect(err).NotTo(HaveOccurred())
	Expect(*dep.Spec.Replicas).To(Equal(int32(1)))

	By("Creating a service for the interconnect default listeners")
	svc, err := f.GetService("edge-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the owner reference for the service")
	Expect(svc.OwnerReferences[0].APIVersion).To(Equal(framework.GVR))
	Expect(svc.OwnerReferences[0].Name).To(Equal("edge-interconnect"))
	Expect(*svc.OwnerReferences[0].Controller).To(Equal(true))

	By("Setting up default listener on qdr instances")
	pods, err := f.ListPodsForDeployment(dep)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods.Items)).To(Equal(1))
	for _, pod := range pods.Items {
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "edge-interconnect", "Router started in Edge mode", time.Second*5)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForRegexpInLog(f.Namespace, pod.Name, "edge-interconnect", `Version:.*1\.8\.0`, time.Second*5)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "edge-interconnect", "Configured Listener: 0.0.0.0:5672 proto=any, role=normal", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "edge-interconnect", "Configured Listener: 0.0.0.0:8080 proto=any, role=normal, http", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "edge-interconnect", "Configured Listener: :8888 proto=any, role=normal, http", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "edge-interconnect", "Configured Listener: 0.0.0.0:55672 proto=any, role=inter-router", time.Second*1)
		Expect(err).To(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "edge-interconnect", "Configured Listener: 0.0.0.0:45672 proto=any, role=edge", time.Second*1)
		Expect(err).To(HaveOccurred())
	}
}
