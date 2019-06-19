package e2e_test

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	. "github.com/interconnectedcloud/qdr-operator/test/e2e"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/client"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/validations"
	"github.com/onsi/ginkgo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// createBasicCR Creates the Basic CR to be validated
func createBasicCR() *v1alpha1.Interconnect {

	// create interconnect custom resource
	// initialize the CR
	basicCR := &v1alpha1.Interconnect{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Interconnect",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "basic-interconnect",
			Namespace: TestSuite.GetNamespace(),
		},
		Spec: v1alpha1.InterconnectSpec{
			DeploymentPlan: v1alpha1.DeploymentPlanType{
				Size:      3,
				Image:     "quay.io/interconnectedcloud/qdrouterd:1.6.0",
				Role:      "interior",
				Placement: "Any",
			},
		},
	}

	return basicCR
}

// TestBasicDeployment Defines the Specs to test a Basic CR using the qdr-operator
func TestBasicDeployment(t *testing.T) {

	var basicCR *v1alpha1.Interconnect = createBasicCR()

	var _ = ginkgo.Describe("Basic", func() {

		ginkgo.Context("Deploy basic-interconnect", func() {
			ginkgo.Describe("Create", func() {
				ginkgo.It("With 3 replicas", func() {
					client.DeployAndValidate(basicCR, t)
				})
			})

			ginkgo.Describe("Update", func() {
				ginkgo.It("DeploymentPlan Size to 4", func() {
					basicCR.Spec.DeploymentPlan.Size = 4
					client.UpdateDeploymentPlanSize(basicCR, t)
				})

				ginkgo.It("DeploymentPlan Size to 3", func() {
					basicCR.Spec.DeploymentPlan.Size = 3
					client.UpdateDeploymentPlanSize(basicCR, t)
				})
			})
		})

		ginkgo.Context("Validate", func() {
			ginkgo.Describe("Entities", func() {
				ginkgo.It("Roles", func() {
					validations.ValidateRoles(basicCR)
				})

				ginkgo.It("RoleBindings", func() {
					validations.ValidateRoleBindings(basicCR)
				})

				ginkgo.It("ServiceAccounts", func() {
					validations.ValidateServiceAccounts(basicCR)
				})

				ginkgo.It("Deployments", func() {
					validations.ValidateDeployment(basicCR)
				})

				ginkgo.It("Pods", func() {
					validations.ValidatePods(basicCR, "3")
				})

				ginkgo.It("Services", func() {
					validations.ValidateService(basicCR, []int{5672, 8080, 55672, 45672})
				})
			})
		})
	})
}
