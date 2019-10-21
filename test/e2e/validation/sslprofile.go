package validation

import (
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/onsi/gomega"
	"k8s.io/api/core/v1"
)

// ValidateDefaultSslProfiles asserts that the default sslProfile entities have
// been defined, based on given Interconnect's role.
func ValidateDefaultSslProfiles(ic *v1alpha1.Interconnect, f *framework.Framework, pods []v1.Pod) {

	var expectedSslProfiles = 1
	var isInterior = ic.Spec.DeploymentPlan.Role == v1alpha1.RouterRoleInterior

	// Interior routers have an extra sslProfile for the inter-router listener
	if isInterior {
		expectedSslProfiles++
	}

	// Iterate through the pods to ensure sslProfiles are defined
	for _, pod := range pods {
		var sslProfilesFound = 0

		// Retrieving sslProfile entities from router
		sslProfiles, err := qdrmanagement.QdmanageQuery(f, pod.Name, entities.SslProfile{}, nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Verify expected sslProfiles are defined
		for _, entity := range sslProfiles {
			sslProfile := entity.(entities.SslProfile)
			switch sslProfile.Name {
			case "inter-router":
				ValidateEntityValues(sslProfile, map[string]interface{}{
					"CaCertFile": fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s-%s-credentials/ca.crt", sslProfile.Name, ic.Name, sslProfile.Name),
				})
				fallthrough
			case "default":
				ValidateEntityValues(sslProfile, map[string]interface{}{
					"CertFile":       fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s-%s-credentials/tls.crt", sslProfile.Name, ic.Name, sslProfile.Name),
					"PrivateKeyFile": fmt.Sprintf("/etc/qpid-dispatch-certs/%s/%s-%s-credentials/tls.key", sslProfile.Name, ic.Name, sslProfile.Name),
				})
				sslProfilesFound++
			}
		}

		// Assert default sslProfiles have been found
		gomega.Expect(expectedSslProfiles).To(gomega.Equal(sslProfilesFound))
	}

}
