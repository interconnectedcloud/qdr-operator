package validations

import (
	"context"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e"
	"github.com/onsi/ginkgo"
	"k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// ValidateDeployment ensures that there is a client with same name as
// the CR name, number of replicas match what is defined in the DeploymentPlan.Size
// and the most recent condition is "Available".
func ValidateDeployment(cr *v1alpha1.Interconnect) {
	cliListOptions := client.ListOptions{
		Namespace: cr.ObjectMeta.Namespace,
	}
	deploymentList := v1.DeploymentList{
		TypeMeta: v12.TypeMeta{
			Kind:       "Deployment",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
	}
	err := e2e.TestSuite.F.Client.List(context.TODO(), &cliListOptions, &deploymentList)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	if len(deploymentList.Items) == 0 {
		ginkgo.Fail(fmt.Sprintf("No deployments found"))
	}
	found := false
	for _, deployment := range deploymentList.Items {

		if deployment.Name != cr.ObjectMeta.Name {
			continue
		}

		found = true

		// Validate replicas
		if deployment.Status.Replicas != cr.Spec.DeploymentPlan.Size {
			ginkgo.Fail(fmt.Sprintf("Invalid replica count. Expected: %v. Found: %v",
				cr.Spec.DeploymentPlan.Size, deployment.Status.Replicas))
		}

		curCondition := deployment.Status.Conditions[0]
		if !strings.EqualFold(string(curCondition.Type), "Available") ||
			!strings.EqualFold(string(curCondition.Status), "True") {
			ginkgo.Fail(fmt.Sprintf("Expected condition/status: Available/True. Got: %v/%v", curCondition.Type, curCondition.Status))
		}
	}
	if !found {
		ginkgo.Fail(fmt.Sprintf("Deployment name not found: %v", cr.ObjectMeta.Name))
	}
}
