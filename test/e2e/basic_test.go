package e2e

import (
	"context"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/validation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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

// testInteriorDefaults creates a simple Interior router using a minimal configuration
// and it asserts that all default elements and artifacts have been defined.
func testInteriorDefaults(f *framework.Framework) {
	var (
		name        = "interior-interconnect"
		defaultSize = 3
		version     = "1.9.0"
	)

	By("Creating a default interior interconnect")
	ei, err := f.CreateInterconnect(f.Namespace, int32(defaultSize), func(ei *v1alpha1.Interconnect) {
		ei.Name = name
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		ei, err := f.GetInterconnect(ei.Name)
		// If not found, no need to delete
		if err != nil && strings.Contains(err.Error(), "not found") {
			return
		}
		Expect(err).NotTo(HaveOccurred())
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

	By("Waiting until full interconnect initial qdr entities")
	ctx, fn = context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx, f, ei)
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

	By("Verifying default listeners have been defined")
	validation.ValidateDefaultListeners(ei, f, pods)

	By("Verifying default connectors have been defined")
	validation.ValidateDefaultConnectors(ei, f, pods)

	By("Verifying default addresses have been defined")
	validation.ValidateDefaultAddresses(ei, f, pods)

	if f.CertManagerPresent {
		By("Verifying expected sslProfiles have been defined")
		validation.ValidateDefaultSslProfiles(ei, f, pods)

		By("Automatically generating credentials")
		_, err = f.GetSecret("interior-interconnect-default-credentials")
		Expect(err).NotTo(HaveOccurred())
		_, err = f.GetSecret("interior-interconnect-inter-router-credentials")
		Expect(err).NotTo(HaveOccurred())
		_, err = f.GetSecret("interior-interconnect-inter-router-ca")
		Expect(err).NotTo(HaveOccurred())
	}

	// Verify service has been defined
	By("Validating service has been defined for Interconnect instance")
	service, err := f.KubeClient.CoreV1().Services(f.Namespace).Get(ei.Name, v1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	Expect(service).NotTo(BeNil())

	// ingress or routes / port <icname-[exposed port #]>
	if !f.IsOpenShift() {
		By("Validating Ingresses have been defined for exposed listeners")
		ingress, err := f.KubeClient.ExtensionsV1beta1().Ingresses(f.Namespace).Get(ei.Name+"-8080", v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(ingress).NotTo(BeNil())

		// Verify exposed inter router TLS listener
		if f.CertManagerPresent {
			ingress, err = f.KubeClient.ExtensionsV1beta1().Ingresses(f.Namespace).Get(ei.Name+"-55671", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(ingress).NotTo(BeNil())
		}
	} else {
		By("Validating Routes have been defined for exposed listeners")
		route, err := f.OcpClient.RoutesClient.RouteV1().Routes(f.Namespace).Get(ei.Name+"-8080", v1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(route).NotTo(BeNil())

		// Verify exposed inter router TLS listener
		if f.CertManagerPresent {
			route, err := f.OcpClient.RoutesClient.RouteV1().Routes(f.Namespace).Get(ei.Name+"-55671", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(route).NotTo(BeNil())
		}
	}

	// Delete and verify related resources were deleted
	By("Deleting the Interconnect instance")
	err = f.DeleteInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	// Waiting till Interconnect is deleted
	err = framework.WaitForDeploymentDeleted(f.KubeClient, f.Namespace, ei.Name, framework.RetryInterval, framework.Timeout)
	Expect(err).NotTo(HaveOccurred())

	// Verify service has been defined
	By("Validating service has been deleted")
	service, err = f.KubeClient.CoreV1().Services(f.Namespace).Get(ei.Name, v1.GetOptions{})
	Expect(err).To(HaveOccurred())

	// ingress or routes / port <icname-[exposed port #]>
	if !f.IsOpenShift() {
		By("Validating Ingresses have been deleted")
		_, err := f.KubeClient.ExtensionsV1beta1().Ingresses(f.Namespace).Get(ei.Name+"-8080", v1.GetOptions{})
		Expect(err).To(HaveOccurred())

		// Verify exposed inter router TLS listener
		if f.CertManagerPresent {
			_, err = f.KubeClient.ExtensionsV1beta1().Ingresses(f.Namespace).Get(ei.Name+"-55671", v1.GetOptions{})
			Expect(err).To(HaveOccurred())
		}
	} else {
		By("Validating Routes have been deleted")
		_, err := f.OcpClient.RoutesClient.RouteV1().Routes(f.Namespace).Get(ei.Name+"-8080", v1.GetOptions{})
		Expect(err).To(HaveOccurred())

		// Verify exposed inter router TLS listener
		if f.CertManagerPresent {
			_, err := f.OcpClient.RoutesClient.RouteV1().Routes(f.Namespace).Get(ei.Name+"-55671", v1.GetOptions{})
			Expect(err).To(HaveOccurred())
		}
	}

}

// testEdgeDefaults defines a minimal edge Interconnect instance and validates
// that all default configuration and resources have been defined.
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
	validation.ValidateDefaultListeners(ei, f, pods)

	By("Verifying default addresses have been defined")
	validation.ValidateDefaultAddresses(ei, f, pods)

	if f.CertManagerPresent {
		By("Verifying expected sslProfiles have been defined")
		validation.ValidateDefaultSslProfiles(ei, f, pods)

		By("Automatically generating credentials")
		_, err = f.GetSecret("edge-interconnect-default-credentials")
		Expect(err).NotTo(HaveOccurred())
	}
}
