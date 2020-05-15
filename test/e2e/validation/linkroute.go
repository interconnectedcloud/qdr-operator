package validation

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

// LinkRouteMapByPrefixPattern represents a map that contains a map
// of string keys that can be either a prefix or a pattern.
type LinkRouteMapByPrefixPattern map[string]map[string]interface{}

// getLinkRouteModel returns the map with the LinkRoute entity model
// if lrMap contains a matching prefix or pattern and a bool value
// that is true if key was found or false otherwise.
func getLinkRouteModel(lrMap LinkRouteMapByPrefixPattern, linkRoute entities.LinkRoute) (map[string]interface{}, bool) {
	emptyModel := map[string]interface{}{}
	lrModel, found := lrMap[linkRoute.Prefix]
	if !found {
		lrModel, found = lrMap[linkRoute.Pattern]
		if !found {
			return emptyModel, false
		}
	}
	return lrModel, true
}

// ValidateSpecLinkRoute asserts that the linkRoute models provided through the lrMap
// are present across all pods from the given ic instance.
func ValidateSpecLinkRoute(ic *v1alpha1.Interconnect, f *framework.Framework, lrMap LinkRouteMapByPrefixPattern) {
	// Retrieving latest Interconnect
	icNew, err := f.GetInterconnect(ic.Name)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Retrieve pod list
	pods, err := f.GetInterconnectPods(icNew)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Iterate through all pods and assert that linkRoutes are available across all instances
	for _, pod := range pods {
		// Wait for pod to be Running
		_, err := f.WaitForPodStatus(pod.Name, v1.PodRunning, framework.Timeout, framework.RetryInterval)
		gomega.Expect(err).To(gomega.BeNil())

		// Same amount of linkRoutes from lrMap are expected to be found
		lrFound := 0

		// Retrieve linkRoutes
		linkRoutes, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.LinkRoute{}, nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Loop through returned linkRoutes
		for _, e := range linkRoutes {
			linkRoute := e.(entities.LinkRoute)
			lrModel, found := getLinkRouteModel(lrMap, linkRoute)
			if !found {
				continue
			}
			// Validating linkRoute that exists on lrMap
			ValidateEntityValues(linkRoute, lrModel)
			lrFound++
		}

		// Assert that all linkRoutes from lrMap have been found
		gomega.Expect(lrFound).To(gomega.Equal(len(lrMap)))
	}
}
