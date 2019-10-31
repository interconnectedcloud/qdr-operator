package e2e

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/validation"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[spec_links] Link Route and Auto Link manipulation tests", func() {

	var (
		icName = "spec-links"
		size   = 3
	)

	// Framework instance to be used across test specs
	f := framework.NewFramework(icName, nil)

	//
	// Validating manipulation of link routes and auto links
	//
	ginkgo.It("Defines link routes and auto links", func() {
		// Create the Interconnect resource with linkRoutes and autoLinks
		ic, err := f.CreateInterconnect(f.Namespace, int32(size), specLinkRoutes, specAutoLinks)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Wait for the deployment to finish
		err = framework.WaitForDeployment(f.KubeClient, f.Namespace, ic.Name, size, framework.RetryInterval, framework.Timeout)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Validating LinkRoute entities
		ginkgo.By("Validating defined LinkRoutes")
		validation.ValidateSpecLinkRoute(ic, f, validation.LinkRouteMapByPrefixPattern{
			"linkRoutePrefixInConnection": {
				"Prefix":            "linkRoutePrefixInConnection",
				"Direction":         common.DirectionTypeIn,
				"Connection":        "connection",
				"AddExternalPrefix": "addExternalPrefix",
				"DelExternalPrefix": "delExternalPrefix",
				"OperStatus":        entities.LinkRouteOperStatusInactive,
			},
			"linkRoutePrefixOutConnection": {
				"Prefix":            "linkRoutePrefixOutConnection",
				"Direction":         common.DirectionTypeOut,
				"Connection":        "connection",
				"AddExternalPrefix": "addExternalPrefix",
				"DelExternalPrefix": "delExternalPrefix",
				"OperStatus":        entities.LinkRouteOperStatusInactive,
			},
			"linkRoutePatternInContainerId": {
				"Pattern":           "linkRoutePatternInContainerId",
				"Direction":         common.DirectionTypeIn,
				"ContainerId":       "containerId",
				"AddExternalPrefix": "addExternalPrefix",
				"DelExternalPrefix": "delExternalPrefix",
				"OperStatus":        entities.LinkRouteOperStatusInactive,
			},
			"linkRoutePatternOutContainerId": {
				"Pattern":           "linkRoutePatternOutContainerId",
				"Direction":         common.DirectionTypeOut,
				"ContainerId":       "containerId",
				"AddExternalPrefix": "addExternalPrefix",
				"DelExternalPrefix": "delExternalPrefix",
				"OperStatus":        entities.LinkRouteOperStatusInactive,
			},
		})

		// Validating AutoLink entities
		ginkgo.By("Validating defined AutoLinks")
		validation.ValidateSpecAutoLink(ic, f, validation.AutoLinkMapByAddress{
			"autoLinkConnection": {
				"Address":         "autoLinkConnection",
				"Direction":       common.DirectionTypeIn,
				"Phase":           1,
				"Connection":      "connection",
				"ExternalAddress": "externalAddress",
				"Fallback":        true,
				"OperStatus":      entities.AutoLinkOperStatusInactive,
			},
			"autoLinkContainerId": {
				"Address":         "autoLinkContainerId",
				"Direction":       common.DirectionTypeOut,
				"Phase":           0,
				"ContainerId":     "containerId",
				"ExternalAddress": "externalAddress",
				"Fallback":        true,
				"OperStatus":      entities.AutoLinkOperStatusInactive,
			},
		})

	})
})

// specLinkRoutes defines a static list of LinkRoutes in
// the provided Interconnect instance.
func specLinkRoutes(interconnect *v1alpha1.Interconnect) {
	interconnect.Spec.LinkRoutes = []v1alpha1.LinkRoute{
		{
			Prefix:            "linkRoutePrefixInConnection",
			Direction:         "in",
			Connection:        "connection",
			AddExternalPrefix: "addExternalPrefix",
			DelExternalPrefix: "delExternalPrefix",
		},
		{
			Prefix:            "linkRoutePrefixOutConnection",
			Direction:         "out",
			Connection:        "connection",
			AddExternalPrefix: "addExternalPrefix",
			DelExternalPrefix: "delExternalPrefix",
		},
		{
			Pattern:           "linkRoutePatternInContainerId",
			Direction:         "in",
			ContainerId:       "containerId",
			AddExternalPrefix: "addExternalPrefix",
			DelExternalPrefix: "delExternalPrefix",
		},
		{
			Pattern:           "linkRoutePatternOutContainerId",
			Direction:         "out",
			ContainerId:       "containerId",
			AddExternalPrefix: "addExternalPrefix",
			DelExternalPrefix: "delExternalPrefix",
		},
	}
}

// specAutoLinks defines a list of AutoLinks in the
// provided Interconnect instance
func specAutoLinks(interconnect *v1alpha1.Interconnect) {
	phaseIn := int32(1)
	phaseOut := int32(0)

	interconnect.Spec.AutoLinks = []v1alpha1.AutoLink{
		{
			Address:         "autoLinkConnection",
			Direction:       "in",
			Connection:      "connection",
			ExternalAddress: "externalAddress",
			Phase:           &phaseIn,
			Fallback:        true,
		},
		{
			Address:         "autoLinkContainerId",
			Direction:       "out",
			ContainerId:     "containerId",
			ExternalAddress: "externalAddress",
			Phase:           &phaseOut,
			Fallback:        true,
		},
	}
}
