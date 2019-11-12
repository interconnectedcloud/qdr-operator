package e2e

import (
	"context"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/validation"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[spec_connector] Connector manipulation tests", func() {

	var (
		icName = "spec-connector"
		size   = 3
	)

	// Framework instance to be used across test specs
	f := framework.NewFramework(icName, nil)

	//
	// Validating manipulation of normal, inter-router and edge connectors
	//
	ginkgo.It("Defines connectors without SSL Profile", func() {
		ginkgo.By("creating an Interconnect with connectors")
		// Create a new Interconnect with the corresponding AMQP connectors
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), specConnectorNoSSL)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Wait till Interconnect up and running
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait till mesh is formed
		ginkgo.By("Waiting until full interconnect initial qdr entities")
		ctx, fn := context.WithTimeout(context.Background(), framework.Timeout)
		defer fn()
		err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx, f, ic)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Retrieve current Interconnect
		ic, err = f.GetInterconnect(ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Validating the defined connectors are present
		validation.ValidateSpecConnector(ic, f, validation.ConnectorMapByPort{
			"5672": {
				"Name": "normal-amqp",
				"Host": "1.1.1.1",
				"Role": common.RoleNormal,
				"Port": "5672",
			},
			"15672": {
				"Name": "inter-amqp",
				"Host": "1.1.1.1",
				"Role": common.RoleInterRouter,
				"Port": "15672",
			},
			"25672": {
				"Name": "edge-amqp",
				"Host": "1.1.1.1",
				"Role": common.RoleEdge,
				"Port": "25672",
			},
		})
	})

	ginkgo.It("Defines connectors with SSL Profile - cert-manager installed", func() {
		if !f.CertManagerPresent {
			ginkgo.Skip("No cert-manager installed")
		}

		ginkgo.By("creating an Interconnect with connectors")
		// Create a new Interconnect with the corresponding AMQPS connectors
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), specConnectorSSL, specConnectorNormalSslProfile)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Wait for deployment
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait till Interconnect up and running
		ginkgo.By("Waiting until full interconnect initial qdr entities")
		ctx, fn := context.WithTimeout(context.Background(), framework.Timeout)
		defer fn()
		err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx, f, ic)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Retrieve current Interconnect
		ic, err = f.GetInterconnect(ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Validating the defined connectors are present
		ginkgo.By("Validating defined connectors")
		validation.ValidateSpecConnector(ic, f, validation.ConnectorMapByPort{
			"5671": {
				"Name":       "normal-amqps",
				"Host":       "1.1.1.1",
				"Role":       common.RoleNormal,
				"SslProfile": "amqps",
				"Port":       "5671",
			},
			"15671": {
				"Name":       "inter-amqps",
				"Host":       "1.1.1.1",
				"Role":       common.RoleInterRouter,
				"SslProfile": "amqps",
				"Port":       "15671",
			},
			"25671": {
				"Name":       "edge-amqps",
				"Host":       "1.1.1.1",
				"Role":       common.RoleEdge,
				"SslProfile": "amqps",
				"Port":       "25671",
			},
		})

		// Validating defined SSL Profiles
		ginkgo.By("Validating SSL Profiles")
		validation.ValidateSslProfileModels(ic, f, validation.SslProfileMapByName{
			"amqps": {
				"Name":           "amqps",
				"CaCertFile":     fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s/ca.crt", "amqps", "amqps-crt"), // amqps-crt because of mutual auth (otherwise would be amqps-ca
				"CertFile":       fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s/tls.crt", "amqps", "amqps-crt"),
				"PrivateKeyFile": fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s/tls.key", "amqps", "amqps-crt"),
			},
		})

		// Verify amqps-ca Issuer exists (as MutualAuth: true)
		ginkgo.By("Validating Issuers")
		issuer, err := f.GetResource(framework.Issuers, "amqps-ca")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(issuer).NotTo(gomega.BeNil())
		gomega.Expect(issuer.GetKind()).To(gomega.Equal("Issuer"))

		// Verify certificates have been generated
		ginkgo.By("Validating Certificates")
		expectedCertificates := []string{"amqps-ca", "amqps-crt"}
		for _, certName := range expectedCertificates {
			cert, err := f.GetResource(framework.Certificates, certName)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(cert).NotTo(gomega.BeNil())
			gomega.Expect(cert.GetKind()).To(gomega.Equal("Certificate"))
		}

		// Deleting the Interconnect instance and validating resources
		err = f.DeleteInterconnect(ic)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Waiting till Interconnect is deleted
		ctx, fn = context.WithTimeout(context.Background(), framework.Timeout)
		defer fn()
		err = framework.WaitForDeploymentDeleted(ctx, f.KubeClient, f.Namespace, ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Verify amqps-ca Issuer has been removed
		ginkgo.By("Validating Issuers removed")
		issuer, err = f.GetResource(framework.Issuers, "amqps-ca")
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(issuer).To(gomega.BeNil())

		// Verify certificates have been removed
		ginkgo.By("Validating Certificates removed")
		for _, certName := range expectedCertificates {
			cert, err := f.GetResource(framework.Certificates, certName)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(cert).To(gomega.BeNil())
		}

	})
})

func specConnectorNoSSL(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.Connectors = []v1alpha1.Connector{
		{
			Name: "normal-amqp",
			Host: "1.1.1.1",
			Port: int32(5672),
		},
	}
	interconnect.Spec.InterRouterConnectors = []v1alpha1.Connector{
		{
			Name: "inter-amqp",
			Host: "1.1.1.1",
			Port: int32(15672),
		},
	}
	interconnect.Spec.EdgeConnectors = []v1alpha1.Connector{
		{
			Name: "edge-amqp",
			Host: "1.1.1.1",
			Port: int32(25672),
		},
	}
}

func specConnectorSSL(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.Connectors = []v1alpha1.Connector{
		{
			Name:       "normal-amqps",
			Host:       "1.1.1.1",
			Port:       int32(5671),
			SslProfile: "amqps",
		},
	}
	interconnect.Spec.InterRouterConnectors = []v1alpha1.Connector{
		{
			Name:       "inter-amqps",
			Host:       "1.1.1.1",
			Port:       int32(15671),
			SslProfile: "amqps",
		},
	}
	interconnect.Spec.EdgeConnectors = []v1alpha1.Connector{
		{
			Name:       "edge-amqps",
			Host:       "1.1.1.1",
			Port:       int32(25671),
			SslProfile: "amqps",
		},
	}
}

func specConnectorNormalSslProfile(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.SslProfiles = []v1alpha1.SslProfile{
		{
			Name:                "amqps",
			Credentials:         "amqps-crt",
			CaCert:              "amqps-ca",
			GenerateCredentials: true,
			GenerateCaCert:      true,
			MutualAuth:          true,
		},
	}
}
