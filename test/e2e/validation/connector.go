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

// ConnectorMapByPort represents a map that contains ports as keys and
// keys/values that represents a Connector entity (for comparison).
type ConnectorMapByPort map[string]map[string]interface{}

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

// ValidateSpecConnector asserts that the connector models provided through the cMap
// are present across all pods from the given ic instance.
func ValidateSpecConnector(ic *v1alpha1.Interconnect, f *framework.Framework, cMap ConnectorMapByPort) {
	// Retrieve fresh version of given Interconnect instance
	icNew, err := f.GetInterconnect(ic.Name)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Iterate through all pods and assert that connectors are available across all instances
	for _, pod := range icNew.Status.PodNames {
		// Same amount of connectors from cMap are expected to be found
		cFound := 0

		// Retrieve connectors
		connectors, err := qdrmanagement.QdmanageQuery(f, pod, entities.Connector{}, nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Loop through returned connectors
		for _, e := range connectors {
			connector := e.(entities.Connector)
			cModel, found := cMap[connector.Port]
			if !found {
				continue
			}
			// Validating that connector exists on cMap
			ValidateEntityValues(connector, cModel)
			cFound++
		}

		// Assert that all connectors from cMap have been found
		gomega.Expect(cFound).To(gomega.Equal(len(cMap)))
	}
}
