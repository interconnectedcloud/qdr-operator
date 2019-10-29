package validation

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"github.com/onsi/gomega"
	"k8s.io/api/core/v1"
)

// ValidateDefaultConnectors asserts that the inter-router connectors are defined
// in [deployment plan size - 1] routers (as the initial pod only provides listeners).
// It returns number of connectors found.
func ValidateDefaultConnectors(interconnect *v1alpha1.Interconnect, f *framework.Framework, pods []v1.Pod) {

	totalConnectors := 0
	expConnectors := 0

	// Expected number of connectors defined by sum of all numbers from 1 to size - 1
	for i := int(interconnect.Spec.DeploymentPlan.Size) - 1; i > 0; i-- {
		expConnectors += i
	}

	// Iterate through pods
	for _, pod := range pods {
		// Retrieve connectors
		connectors, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.Connector{}, nil)
		gomega.Expect(err).To(gomega.BeNil())

		// Common connector properties
		expectedPort := "55672"
		expectedSslProfile := ""
		if f.CertManagerPresent {
			expectedPort = "55671"
			expectedSslProfile = "inter-router"
		}

		props := map[string]interface{}{
			"Role":       common.RoleInterRouter,
			"Port":       expectedPort,
			"SslProfile": expectedSslProfile,
		}

		// Validate connectors
		if len(connectors) > 0 {
			for _, entity := range connectors {
				connector := entity.(entities.Connector)
				gomega.Expect(connector.Host).NotTo(gomega.BeEmpty())
				ValidateEntityValues(connector, props)
			}
		}

		totalConnectors += len(connectors)
	}

	// Validate number of connectors across pods
	gomega.Expect(expConnectors).To(gomega.Equal(totalConnectors))

}
