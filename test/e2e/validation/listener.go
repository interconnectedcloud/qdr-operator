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

// ListenerMapByPort represents a map that contains
type ListenerMapByPort map[string]map[string]interface{}

// ValidateDefaultListeners ensures that the default listeners (if no others specified)
// have been created
func ValidateDefaultListeners(ic *v1alpha1.Interconnect, f *framework.Framework, pods []v1.Pod) {
	var expLs int = 3
	var expIntLs int = 2

	// Expected inter-router listener port might change if CertManager is present
	var interRouterPort = "55672"
	var interRouterSslProfile = ""
	if f.CertManagerPresent {
		expLs++
		interRouterPort = "55671"
		interRouterSslProfile = "inter-router"
	}

	for _, pod := range pods {
		var lsFound int = 0
		var intLsFound int = 0

		listeners, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.Listener{}, nil)
		gomega.Expect(err).To(gomega.BeNil())

		// Validate returned listeners
		for _, e := range listeners {
			l := e.(entities.Listener)
			switch l.Port {
			case "5671":
				ValidateEntityValues(l, map[string]interface{}{
					"Port": "5671",
					"Role": common.RoleNormal,
					"Http": false,
				})
				lsFound++
			case "5672":
				ValidateEntityValues(l, map[string]interface{}{
					"Port": "5672",
					"Role": common.RoleNormal,
					"Http": false,
				})
				lsFound++
			case "8080":
				ValidateEntityValues(l, map[string]interface{}{
					"Port":             "8080",
					"Role":             common.RoleNormal,
					"Http":             true,
					"AuthenticatePeer": true,
				})
				lsFound++
			case "8888":
				ValidateEntityValues(l, map[string]interface{}{
					"Port":        "8888",
					"Role":        common.RoleNormal,
					"Http":        true,
					"Healthz":     true,
					"Metrics":     true,
					"Websockets":  false,
					"HttpRootDir": "invalid",
				})
				lsFound++
			}
		}

		// Expect default listener count to match
		gomega.Expect(expLs).To(gomega.Equal(lsFound))

		//
		// Interior only
		//
		if ic.Spec.DeploymentPlan.Role != v1alpha1.RouterRoleInterior {
			return
		}

		// Validate interior listeners
		for _, e := range listeners {
			l := e.(entities.Listener)
			switch l.Port {
			// inter-router listener
			case interRouterPort:
				ValidateEntityValues(l, map[string]interface{}{
					"Port":       interRouterPort,
					"Role":       common.RoleInterRouter,
					"SslProfile": interRouterSslProfile,
				})
				intLsFound++
			// edge listener
			case "45672":
				ValidateEntityValues(l, map[string]interface{}{
					"Port": "45672",
					"Role": common.RoleEdge,
				})
				intLsFound++
			}
		}

		// Validate all default interior listeners are present
		gomega.Expect(expIntLs).To(gomega.Equal(intLsFound))
	}
}

// ValidateSpecListener asserts that the listener models provided through the lsMap
// are present across all pods from the given ic instance.
func ValidateSpecListener(ic *v1alpha1.Interconnect, f *framework.Framework, lsMap ListenerMapByPort) {
	// Retrieve fresh version of given Interconnect instance
	icNew, err := f.GetInterconnect(ic.Name)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Iterate through all pods and assert that listeners are available across all instances
	for _, pod := range icNew.Status.PodNames {
		// Same amount of listeners from lsMap are expected to be found
		lsFound := 0

		// Retrieve listeners
		listeners, err := qdrmanagement.QdmanageQuery(f, pod, entities.Listener{}, nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Loop through returned listeners
		for _, e := range listeners {
			listener := e.(entities.Listener)
			lsModel, found := lsMap[listener.Port]
			if !found {
				continue
			}
			// Validating listener that exists on lsMap
			ValidateEntityValues(listener, lsModel)
			lsFound++
		}

		// Assert that all listeners from lsMap have been found
		gomega.Expect(lsFound).To(gomega.Equal(len(lsMap)))
	}
}
