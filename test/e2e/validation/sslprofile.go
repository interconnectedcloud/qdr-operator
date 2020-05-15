package validation

import (
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	"github.com/onsi/gomega"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SslProfileMapByName represents a map indexed by sslProfile Name storing
// another map with the property names and respective values for the SslProfile entity
// that will be validated.
type SslProfileMapByName map[string]map[string]interface{}

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

// ValidateSslProfileModels retrieves the Interconnect instance and iterates through all
// its pods, querying management API for sslProfiles. Next it ensure that all sslProfile
// definitions fro the sslProfMap are defined on each pod.
func ValidateSslProfileModels(ic *v1alpha1.Interconnect, f *framework.Framework, sslProfMap SslProfileMapByName) {
	var podNames []string

	// Retrieve lastest version of given Interconnect resource
	ic, err := f.GetInterconnect(ic.Name)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// Validate IC instance
	gomega.Expect(ic).NotTo(gomega.BeNil())

	pods, err := f.KubeClient.CoreV1().Pods(ic.Namespace).List(metav1.ListOptions{LabelSelector: "application=" + ic.Name + ",interconnect_cr=" + ic.Name})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Expect(len(pods.Items)).To(gomega.BeNumerically(">", 0))

	for _, pod := range pods.Items {
		if pod.GetObjectMeta().GetDeletionTimestamp() == nil {
			podNames = append(podNames, pod.Name)
		}
	}

	for _, pod := range podNames {
		sslProfFound := 0

		sslProfiles, err := qdrmanagement.QdmanageQuery(f, pod, entities.SslProfile{}, nil)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		for _, e := range sslProfiles {
			sslProfile := e.(entities.SslProfile)
			model, found := sslProfMap[sslProfile.Name]
			if !found {
				continue
			}

			ValidateEntityValues(sslProfile, model)
			// Validating the matching sslProfile
			sslProfFound++
		}

		// Expect all sslProfiles from map have been validated
		gomega.Expect(sslProfFound).To(gomega.Equal(len(sslProfMap)))
	}
}
