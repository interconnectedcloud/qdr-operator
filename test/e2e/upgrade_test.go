package e2e

import (
	//"time"

	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"

	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[upgrade_test] Interconnect upgrade deployment tests", func() {
	f := framework.NewFramework("basic-interior", nil)

	It("Should be able to upgrade the qdrouterd image for an interior deployment", func() {
		testInteriorImageUpgrade(f)
	})

})

func testInteriorImageUpgrade(f *framework.Framework) {
	By("Creating an interior interconnect with size 3")
	ei, err := f.CreateInterconnect(f.Namespace, 3, func(ei *v1alpha1.Interconnect) {
		ei.Name = "interior-interconnect"
		ei.Spec.DeploymentPlan.Image = "quay.io/interconnectedcloud/qdrouterd:1.8.0"
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

	By("Creating a Deployment resource in the namespace")
	dep, err := f.GetDeployment("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment contains 3 pods")
	pods, err := f.ListPodsForDeployment(dep)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods.Items)).To(Equal(3))

	By("Verifying the Network contains 3 nodes on each pod")
	for _, pod := range pods.Items {
		err := qdrmanagement.WaitForQdrNodesInPod(f, pod, 3, framework.RetryInterval, framework.Timeout)
		Expect(err).NotTo(HaveOccurred())
		nodes, err := qdrmanagement.QdmanageQueryNodes(f, pod.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(nodes)).To(Equal(3))
	}

	By("Upgrading the qdrouterd")
	ei.Spec.DeploymentPlan.Image = "quay.io/interconnectedcloud/qdrouterd:1.9.0"
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment contains 3 pods")
	pods, err = f.ListPodsForDeployment(dep)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods.Items)).To(Equal(3))

	//By("Verifying the pods are running the upgrade image")
	//pods, err = f.ListPodsForDeployment(dep)
	//Expect(err).NotTo(HaveOccurred())
	//Expect(len(pods.Items)).To(Equal(3))
	//for _, pod := range pods.Items {
	//	if pod.GetObjectMeta().GetDeletionTimestamp() == nil {
	//		_, err = framework.LookForRegexpInLog(f.Namespace, pod.Name, "interior-interconnect", `Version:.*1\.9\.0`, time.Second*20)
	//		Expect(err).NotTo(HaveOccurred())
	//	}
	//}

	//By("Verifying each node has 2 inter-router connections")
	//for _, pod := range pods.Items {
	//    if pod.GetObjectMeta().GetDeletionTimestamp() == nil {
	//	// Retrieving inter-router connections 2 on each of the 3 nodes
	//	conns, err := qdrmanagement.ListInterRouterConnectionsForPod(f, pod)
	//	Expect(err).NotTo(HaveOccurred())
	//	Expect(len(conns)).To(Equal(2))
	//    }
	//}

}
