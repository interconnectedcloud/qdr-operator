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
