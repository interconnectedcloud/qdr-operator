package e2e

import (
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	corev1 "k8s.io/api/core/v1"

	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[resize_test] Interconnect resize deployment tests", func() {
	f := framework.NewFramework("basic-interior", nil)

	It("Should be able to resize an interior deployment from 3 to 5", func() {
		testInteriorResizeFrom3To5(f)
	})

	It("Should be able to resize an interior deployment from 5 to 3", func() {
		testInteriorResizeFrom5To3(f)
	})
})

func testInteriorResizeFrom3To5(f *framework.Framework) {
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

	By("Verifying the deployment in the namespace")
	dep, err := f.GetDeployment("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan size")
	Expect(ei.Spec.DeploymentPlan.Size).To(Equal(int32(3)))

	By("Scaling the interior interconnect size to 5")
	ei.Spec.DeploymentPlan.Size = 5
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment has reached 5 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 5, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment contains 5 pods")
	pods, err := f.ListPodsForDeployment(dep)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods.Items)).To(Equal(5))

	By("Verifying the Network contains 5 nodes on each pod")
	for _, pod := range pods.Items {
		err := qdrmanagement.WaitForQdrNodesInPod(f, pod, 5, framework.RetryInterval, framework.Timeout)
		Expect(err).NotTo(HaveOccurred())
		nodes, err := qdrmanagement.QdmanageQueryNodes(f, pod.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(nodes)).To(Equal(5))
	}

	By("Verifying each node has 4 inter-router connections")
	for _, pod := range pods.Items {
		// Retrieving inter-router connections 4 on each of the 5 nodes
		conns, err := qdrmanagement.ListInterRouterConnectionsForPod(f, pod)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(conns)).To(Equal(4))
	}

}

func testInteriorResizeFrom5To3(f *framework.Framework) {
	var activepods []corev1.Pod

	By("Creating an interior interconnect with size 5")
	ei, err := f.CreateInterconnect(f.Namespace, 5, func(ei *v1alpha1.Interconnect) {
		ei.Name = "interior-interconnect"
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating a Deployment with 5 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 5, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment in the namespace")
	dep, err := f.GetDeployment("interior-interconnect")
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan size")
	Expect(ei.Spec.DeploymentPlan.Size).To(Equal(int32(5)))

	By("Scaling the interior interconnect size to 3")
	ei.Spec.DeploymentPlan.Size = 3
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment has reached 3 replicas")
	err = framework.WaitForDeployment(f.KubeClient, f.Namespace, "interior-interconnect", 3, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the Deployment contains 3 pods")
	pods, err := f.ListPodsForDeployment(dep)
	Expect(err).NotTo(HaveOccurred())

	for _, pod := range pods.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() == nil {
			activepods = append(activepods, pod)
		}
	}
	Expect(len(activepods)).To(Equal(3))

	By("Verifying the Network contains 3 nodes on each pod")
	for _, pod := range activepods {
		err := qdrmanagement.WaitForQdrNodesInPod(f, pod, 3, framework.RetryInterval, framework.Timeout)
		Expect(err).NotTo(HaveOccurred())
	}

	By("Verifying each node has 2 inter-router connections")
	for _, pod := range activepods {
		// Retrieving inter-router connections 2 on each of the 3 nodes
		conns, err := qdrmanagement.ListInterRouterConnectionsForPod(f, pod)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(conns)).To(Equal(2))
	}

}
