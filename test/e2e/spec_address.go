package e2e

import (
	"context"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/validation"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[spec_address] Address manipulation tests", func() {
	var (
		icName          = "spec-address"
		size            = 3
		newAddressesMap = map[string]map[string]interface{}{
			"multicastprefix": {
				"Prefix":       "multicastprefix",
				"Distribution": entities.DistributionMulticast,
			},
			"multicastpattern/*": {
				"Pattern":      "multicastpattern/*",
				"Distribution": entities.DistributionMulticast,
			},
			"waypoint": {
				"Prefix":       "waypoint",
				"Distribution": entities.DistributionBalanced,
				"Waypoint":     true,
			},
			"fallback": {
				"Prefix":         "fallback",
				"Distribution":   entities.DistributionBalanced,
				"Waypoint":       true,
				"EnableFallback": true,
			},
			"priority": {
				"Prefix":   "priority",
				"Priority": 5,
			},
			"ingressegress": {
				"Prefix":       "ingressegress",
				"IngressPhase": 5,
				"EgressPhase":  6,
			},
		}
		modifiedAddressesMap = map[string]map[string]interface{}{
			"modified/multicastprefix": {
				"Prefix":       "modified/multicastprefix",
				"Distribution": entities.DistributionMulticast,
				"Priority":     4,
			},
			"modified/multicastpattern/*": {
				"Pattern":      "modified/multicastpattern/*",
				"Distribution": entities.DistributionMulticast,
				"Priority":     4,
			},
			"modified/waypoint": {
				"Prefix":       "modified/waypoint",
				"Distribution": entities.DistributionBalanced,
				"Waypoint":     true,
				"Priority":     4,
			},
			"modified/fallback": {
				"Prefix":         "modified/fallback",
				"Distribution":   entities.DistributionBalanced,
				"Waypoint":       true,
				"EnableFallback": true,
				"Priority":       4,
			},
			"modified/priority": {
				"Prefix":   "modified/priority",
				"Priority": 4,
			},
			"modified/ingressegress": {
				"Prefix":       "modified/ingressegress",
				"IngressPhase": 5,
				"EgressPhase":  6,
				"Priority":     4,
			},
		}
		afterRemovalAddressesMap = map[string]map[string]interface{}{
			"multicastprefix": {
				"Prefix":       "multicastprefix",
				"Distribution": entities.DistributionMulticast,
			},
			"multicastpattern/*": {
				"Pattern":      "multicastpattern/*",
				"Distribution": entities.DistributionMulticast,
			},
			"waypoint": {
				"Prefix":       "waypoint",
				"Distribution": entities.DistributionBalanced,
				"Waypoint":     true,
			},
		}
	)

	// Create the framework instance
	f := framework.NewFramework(icName, nil)

	ginkgo.It("Should be able to define addresses", func() {
		// Deploy and validate deployment
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), func(ic *v1alpha1.Interconnect) {
			ic.Spec.Addresses = getSpecAddresses()
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Wait till interconnect instance is ready
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Validating addresses
		validateAddresses(f, ic, newAddressesMap)
	})

	ginkgo.It("Should be able to modify defined addresses", func() {
		// Deploy and validate deployment
		ginkgo.By("Deploying an Interconnect using pre-defined address list")
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), func(ic *v1alpha1.Interconnect) {
			ic.Spec.Addresses = getSpecAddresses()
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Wait till interconnect instance is ready
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Get the Interconnect instance
		ic, err = f.GetInterconnect(ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Modify Prefix/Pattern and Priority of all addresses
		newAddresses := []v1alpha1.Address{}
		for _, addr := range ic.Spec.Addresses {
			if addr.Pattern != "" {
				addr.Pattern = "modified/" + addr.Pattern
			} else {
				addr.Prefix = "modified/" + addr.Prefix
			}
			priority := int32(4)
			addr.Priority = &priority
			newAddresses = append(newAddresses, addr)
		}
		ic.Spec.Addresses = newAddresses

		// Updating Interconnect resource
		ginkgo.By("Modifying addresses and updating Interconnect instance")
		ic, err = f.UpdateInterconnect(ic)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Waiting for new pods
		ctx, fn := context.WithTimeout(context.Background(), framework.Timeout)
		defer fn()
		err = f.WaitForNewInterconnectPods(ctx, ic, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait till interconnect instance is ready
		ic, err = f.GetInterconnect(ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Waiting for deployment to be ready
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Validating addresses
		ginkgo.By("Asserting that addresses have been updated accordingly")
		validateAddresses(f, ic, modifiedAddressesMap)
	})

	ginkgo.It("Should be able to remove defined addresses", func() {
		// Deploy and validate deployment
		ginkgo.By("Deploying an Interconnect using pre-defined address list")
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), func(ic *v1alpha1.Interconnect) {
			ic.Spec.Addresses = getSpecAddresses()
		})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Wait till interconnect instance is ready
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(ic).NotTo(gomega.BeNil())

		// Get the Interconnect instance
		ic, err = f.GetInterconnect(ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Keeping just the first 3 addresses
		validateAddresses(f, ic, newAddressesMap)

		// Updating the list of addresses
		newAddrList := getSpecAddresses()[:3]
		ic.Spec.Addresses = newAddrList

		// Updating Interconnect resource
		ginkgo.By("Modifying addresses and updating Interconnect instance")
		ic, err = f.UpdateInterconnect(ic)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Waiting for new pods
		ctx, fn := context.WithTimeout(context.Background(), framework.Timeout)
		defer fn()
		err = f.WaitForNewInterconnectPods(ctx, ic, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait till interconnect instance is ready
		ic, err = f.GetInterconnect(ic.Name)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Waiting for deployment to be ready
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Validating addresses
		ginkgo.By("Asserting that addresses have been updated accordingly")
		validateAddresses(f, ic, afterRemovalAddressesMap)
	})
})

// validateAddresses asserts that the Address entities available across all
// pods of the given Interconnect instance match the definitions from addrMap.
func validateAddresses(f *framework.Framework, ic *v1alpha1.Interconnect, addrMap map[string]map[string]interface{}) {
	// Validating defined addresses/
	ic, err := f.GetInterconnect(ic.Name)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Retrieve current pod names
	podNames := ic.Status.PodNames
	for _, pod := range podNames {
		addresses, err := qdrmanagement.QdmanageQuery(f, pod, entities.Address{}, nil)
		gomega.Expect(err).To(gomega.BeNil())
		// Assert number of addresses match expected count
		gomega.Expect(len(addrMap)).To(gomega.Equal(len(addresses)))

		// Validating all addresses
		for _, address := range addresses {
			a := getAddressMapForPrefixPattern(addrMap, address.(entities.Address))
			gomega.Expect(len(addrMap)).To(gomega.BeNumerically(">", 0))
			validation.ValidateEntityValues(address, a)
		}
	}
}

func getSpecAddresses() []v1alpha1.Address {

	priority := int32(5)
	ingress := int32(5)
	egress := int32(6)

	return []v1alpha1.Address{
		{
			Prefix:       "multicastprefix",
			Distribution: "multicast",
		},
		{
			Pattern:      "multicastpattern/*",
			Distribution: "multicast",
		},
		{
			Prefix:       "waypoint",
			Distribution: "balanced",
			Waypoint:     true,
		},
		{
			Prefix:         "fallback",
			Distribution:   "balanced",
			Waypoint:       true,
			EnableFallback: true,
		},
		{
			Prefix:   "priority",
			Priority: &priority,
		},
		{
			Prefix:       "ingressegress",
			IngressPhase: &ingress,
			EgressPhase:  &egress,
		},
	}
}

func getAddressMapForPrefixPattern(m map[string]map[string]interface{}, address entities.Address) map[string]interface{} {
	for k, v := range m {
		if k == address.Prefix || k == address.Pattern {
			return v
		}
	}
	return map[string]interface{}{}
}
