package e2e

import (
	"context"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"time"

	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[basic_test] Interconnect default deployment tests", func() {
	f := framework.NewFramework("basic-interior", nil)

	It("Should be able to create a default interior deployment", func() {
		testInteriorDefaults(f)
	})

	It("Should be able to create a default edge deployment", func() {
		testEdgeDefaults(f)
	})

})

func testInteriorDefaults(f *framework.Framework) {
	var (
		name        = "interior-interconnect"
		defaultSize = 1
		version     = "1.9.0"
	)

	By("Creating a default interior interconnect")
	ei, err := f.CreateInterconnect(f.Namespace, 0, func(ei *v1alpha1.Interconnect) {
		ei.Name = name
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect(name)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan defaults")
	Expect(ei.Name).To(Equal(name))
	Expect(ei.Spec.DeploymentPlan.Role).To(Equal(v1alpha1.RouterRoleInterior))
	Expect(ei.Spec.DeploymentPlan.Placement).To(Equal(v1alpha1.PlacementAny))

	By("Waiting until full interconnect with version")
	ctx, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = f.WaitUntilFullInterconnectWithVersion(ctx, ei, defaultSize, version)
	Expect(err).NotTo(HaveOccurred())

	By("Creating a service for the interconnect default listeners")
	svc, err := f.GetService(name)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the owner reference for the service")
	Expect(svc.OwnerReferences[0].APIVersion).To(Equal(framework.GVR))
	Expect(svc.OwnerReferences[0].Name).To(Equal(name))
	Expect(*svc.OwnerReferences[0].Controller).To(Equal(true))

	By("Setting up default listener on qdr instances")
	pods, err := f.PodsForInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods)).To(Equal(defaultSize))
	for _, pod := range pods {
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:5672 proto=any, role=normal", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:8080 proto=any, role=normal, http", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: :8888 proto=any, role=normal, http", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		if f.CertManagerPresent {
			_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:55671 proto=any, role=inter-router, sslProfile=inter-router", time.Second*1)
		} else {
			_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:55672 proto=any, role=inter-router", time.Second*1)
		}
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, "interior-interconnect", "Configured Listener: 0.0.0.0:45672 proto=any, role=edge", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
	}

	By("Verifying each node has 5 addresses")
	for _, pod := range pods {
		addrs, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.Address{}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(addrs)).To(Equal(5))
	}

	if f.CertManagerPresent {
		By("Automatically generating credentials")
		_, err = f.GetSecret("interior-interconnect-default-credentials")
		Expect(err).NotTo(HaveOccurred())
		_, err = f.GetSecret("interior-interconnect-inter-router-credentials")
		Expect(err).NotTo(HaveOccurred())
		_, err = f.GetSecret("interior-interconnect-inter-router-ca")
		Expect(err).NotTo(HaveOccurred())
	}
}

func testEdgeDefaults(f *framework.Framework) {
	var (
		name        = "edge-interconnect"
		role        = "edge"
		defaultSize = 1
		version     = "1.9.0"
	)

	By("Creating an edge interconnect with default size")
	ei, err := f.CreateInterconnect(f.Namespace, 0, func(ei *v1alpha1.Interconnect) {
		ei.Name = name
		ei.Spec.DeploymentPlan.Role = v1alpha1.RouterRoleType(role)
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect(name)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the deployment plan")
	Expect(ei.Name).To(Equal(name))
	Expect(ei.Spec.DeploymentPlan.Role).To(Equal(v1alpha1.RouterRoleType(role)))
	Expect(ei.Spec.DeploymentPlan.Placement).To(Equal(v1alpha1.PlacementAny))

	By("Waiting until full interconnect with version")
	ctx, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = f.WaitUntilFullInterconnectWithVersion(ctx, ei, defaultSize, version)
	Expect(err).NotTo(HaveOccurred())

	By("Creating a service for the interconnect default listeners")
	svc, err := f.GetService(name)
	Expect(err).NotTo(HaveOccurred())

	By("Verifying the owner reference for the service")
	Expect(svc.OwnerReferences[0].APIVersion).To(Equal(framework.GVR))
	Expect(svc.OwnerReferences[0].Name).To(Equal(name))
	Expect(*svc.OwnerReferences[0].Controller).To(Equal(true))

	By("Setting up default listener on qdr instances")
	pods, err := f.PodsForInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())
	Expect(len(pods)).To(Equal(defaultSize))
	for _, pod := range pods {
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Router started in Edge mode", time.Second*5)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:5672 proto=any, role=normal", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		if f.CertManagerPresent {
			_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:5671 proto=any, role=normal, sslProfile=default", time.Second*1)
			Expect(err).NotTo(HaveOccurred())
		}
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:8080 proto=any, role=normal, http", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: :8888 proto=any, role=normal, http", time.Second*1)
		Expect(err).NotTo(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:55672 proto=any, role=inter-router", time.Second*1)
		Expect(err).To(HaveOccurred())
		_, err = framework.LookForStringInLog(f.Namespace, pod.Name, name, "Configured Listener: 0.0.0.0:45672 proto=any, role=edge", time.Second*1)
		Expect(err).To(HaveOccurred())
	}

	By("Verifying each node has 5 addresses")
	for _, pod := range pods {
		addrs, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.Address{}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(addrs)).To(Equal(5))
	}

	if f.CertManagerPresent {
		By("Automatically generating credentials")
		_, err = f.GetSecret("edge-interconnect-default-credentials")
		Expect(err).NotTo(HaveOccurred())
		_, err = f.GetSecret("edge-interconnect-inter-router-credentials")
		Expect(err).NotTo(HaveOccurred())
		_, err = f.GetSecret("edge-interconnect-inter-router-ca")
		Expect(err).NotTo(HaveOccurred())
	}
}
