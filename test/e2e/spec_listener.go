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

var _ = ginkgo.Describe("[spec_listener] Listener manipulation tests", func() {

	var (
		icName = "spec-listener"
		size   = 3
	)

	// Framework instance to be used across test specs
	f := framework.NewFramework(icName, nil)

	//
	// Validating manipulation of normal, inter-router and edge listeners
	//
	ginkgo.It("Defines listeners", func() {
		ginkgo.By("creating an Interconnect with listeners")
		// Create a new Interconnect with the corresponding AMQP and HTTP listeners
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), specListenerNoSSL)
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

		// Validating the two defined listeners are present
		validation.ValidateSpecListener(ic, f, validation.ListenerMapByPort{
			"5672": {
				"Name": "normal-amqp",
				"Host": "0.0.0.0",
				"Role": common.RoleNormal,
				"Port": "5672",
			},
			"8080": {
				"Name": "normal-http",
				"Host": "0.0.0.0",
				"Role": common.RoleNormal,
				"Port": "8080",
				"Http": true,
			},
			"15672": {
				"Name": "inter-amqp",
				"Host": "0.0.0.0",
				"Role": common.RoleInterRouter,
				"Port": "15672",
			},
			"25672": {
				"Name": "edge-amqp",
				"Host": "0.0.0.0",
				"Role": common.RoleEdge,
				"Port": "25672",
			},
		})
	})

	ginkgo.It("Defines listeners with cert-manager installed", func() {
		if !f.CertManagerPresent {
			ginkgo.Skip("No cert-manager installed")
		}

		ginkgo.By("creating an Interconnect with AMQPS and HTTPS listeners")
		// Create a new Interconnect with the corresponding AMQP and HTTP listeners
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), specListenerSSL, specListenerNormalSslProfile)
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

		// Validating the two defined listeners are present
		ginkgo.By("Validating defined listeners")
		validation.ValidateSpecListener(ic, f, validation.ListenerMapByPort{
			"5672": {
				"Name": "normal-amqp",
				"Host": "0.0.0.0",
				"Role": common.RoleNormal,
				"Port": "5672",
			},
			"5671": {
				"Name":       "normal-amqps",
				"Host":       "0.0.0.0",
				"Role":       common.RoleNormal,
				"SslProfile": "amqps",
				"Port":       "5671",
			},
			"8443": {
				"Name":       "normal-https",
				"Host":       "0.0.0.0",
				"Role":       common.RoleNormal,
				"SslProfile": "https",
				"Port":       "8443",
				"Http":       true,
			},
			"15671": {
				"Name":           "inter-amqps",
				"Host":           "0.0.0.0",
				"Role":           common.RoleInterRouter,
				"SslProfile":     "amqps",
				"SaslMechanisms": "EXTERNAL",
				"Port":           "15671",
			},
			"25671": {
				"Name":       "edge-amqps",
				"Host":       "0.0.0.0",
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
			"https": {
				"Name":           "https",
				"CaCertFile":     fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s/tls.crt", "https", "https-ca"),
				"CertFile":       fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s/tls.crt", "https", "https-crt"),
				"PrivateKeyFile": fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s/tls.key", "https", "https-crt"),
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
		expectedCertificates := []string{"amqps-ca", "amqps-crt", "https-ca", "https-crt"}
		for _, certName := range expectedCertificates {
			cert, err := f.GetResource(framework.Certificates, certName)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(cert).NotTo(gomega.BeNil())
			gomega.Expect(cert.GetKind()).To(gomega.Equal("Certificate"))
		}

	})
})

func specListenerNoSSL(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.Listeners = []v1alpha1.Listener{
		{
			Name:   "normal-amqp",
			Host:   "0.0.0.0",
			Port:   int32(5672),
			Expose: true,
		},
		{
			Name:   "normal-http",
			Host:   "0.0.0.0",
			Port:   int32(8080),
			Http:   true,
			Expose: true,
		},
	}
	interconnect.Spec.InterRouterListeners = []v1alpha1.Listener{
		{
			Name:   "inter-amqp",
			Host:   "0.0.0.0",
			Port:   int32(15672),
			Expose: true,
		},
	}
	interconnect.Spec.EdgeListeners = []v1alpha1.Listener{
		{
			Name:   "edge-amqp",
			Host:   "0.0.0.0",
			Port:   int32(25672),
			Expose: true,
		},
	}
}

func specListenerSSL(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.Listeners = []v1alpha1.Listener{
		{
			Name:   "normal-amqp",
			Host:   "0.0.0.0",
			Port:   int32(5672),
			Expose: true,
		},
		{
			Name:             "normal-amqps",
			Host:             "0.0.0.0",
			Port:             int32(5671),
			SslProfile:       "amqps",
			AuthenticatePeer: true,
			Expose:           true,
		},
		{
			Name:             "normal-https",
			Host:             "0.0.0.0",
			Port:             int32(8443),
			Http:             true,
			SslProfile:       "https",
			AuthenticatePeer: false,
			Expose:           true,
		},
	}
	interconnect.Spec.InterRouterListeners = []v1alpha1.Listener{
		{
			Name:             "inter-amqps",
			Host:             "0.0.0.0",
			Port:             int32(15671),
			SslProfile:       "amqps",
			SaslMechanisms:   "EXTERNAL",
			AuthenticatePeer: true,
			Expose:           true,
		},
	}
	interconnect.Spec.EdgeListeners = []v1alpha1.Listener{
		{
			Name:             "edge-amqps",
			Host:             "0.0.0.0",
			Port:             int32(25671),
			SslProfile:       "amqps",
			SaslMechanisms:   "EXTERNAL",
			AuthenticatePeer: true,
			Expose:           true,
		},
	}
}

func specListenerNormalSslProfile(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.SslProfiles = []v1alpha1.SslProfile{
		{
			Name:                "amqps",
			Credentials:         "amqps-crt",
			CaCert:              "amqps-ca",
			GenerateCredentials: true,
			GenerateCaCert:      true,
			MutualAuth:          true,
		},
		{
			Name:                "https",
			Credentials:         "https-crt",
			CaCert:              "https-ca",
			GenerateCredentials: true,
			GenerateCaCert:      true,
			MutualAuth:          false,
		},
	}
}
