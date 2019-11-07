package validation

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/onsi/gomega"
	"k8s.io/api/core/v1"
)

// ValidateDefaultAddresses verifies that the created addresses match expected ones
func ValidateDefaultAddresses(ic *v1alpha1.Interconnect, f *framework.Framework, pods []v1.Pod) {

	const expectedAddresses = 5

	for _, pod := range pods {
		var defaultAddressesFound = 0

		// Querying addresses on given pod
		addrs, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.Address{}, nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(len(addrs)).To(gomega.Equal(expectedAddresses))

		// Validates all addresses are present and match expected definition
		for _, entity := range addrs {
			addr := entity.(entities.Address)
			switch addr.Prefix {
			case "closest":
				fallthrough
			case "unicast":
				fallthrough
			case "exclusive":
				ValidateEntityValues(addr, map[string]interface{}{
					"Distribution": entities.DistributionClosest,
				})
				defaultAddressesFound++
			case "multicast":
				fallthrough
			case "broadcast":
				ValidateEntityValues(addr, map[string]interface{}{
					"Distribution": entities.DistributionMulticast,
				})
				defaultAddressesFound++
			}
		}

		// Assert default addresses have been found
		gomega.Expect(expectedAddresses).To(gomega.Equal(defaultAddressesFound))
	}

}
