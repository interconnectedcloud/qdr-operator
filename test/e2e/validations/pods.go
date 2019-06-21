package validations

import (
	"context"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/utils"
	"github.com/onsi/ginkgo"
	v13 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

// ValidatePods queries all Pods from current namespace and validate:
// - number of pods matching CR name prefix must match DeploymentPlan.Size
// - each pod must have the following EnvVars:
//   - QDROUTERD_CONF (with a router and at least one listener element)
//   - APPLICATION_NAME (matching CR name)
//   - POD_COUNT (match original replica size)
//
func ValidatePods(cr *v1alpha1.Interconnect, originalReplicaSize string) {
	cliListOptions := client.ListOptions{
		Namespace: cr.ObjectMeta.Namespace,
	}
	podList := v1.PodList{
		TypeMeta: v12.TypeMeta{
			Kind:       "Pod",
			APIVersion: v13.SchemeGroupVersion.String(),
		},
	}
	err := e2e.TestSuite.F.Client.List(context.TODO(), &cliListOptions, &podList)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	if len(podList.Items) == 0 {
		ginkgo.Fail(fmt.Sprintf("No Pods found"))
	}
	count := 0
	expEnvVars := []string{"APPLICATION_NAME", "QDROUTERD_CONF", "POD_COUNT"}
	for _, pod := range podList.Items {
		// If pod name does not match expected prefix or pod has been deleted, skip it
		if !strings.HasPrefix(pod.Name, cr.ObjectMeta.Name+"-") || pod.DeletionTimestamp != nil {
			continue
		}

		count++
		if "Running" != pod.Status.Phase {
			ginkgo.Fail(fmt.Sprintf("Invalid POD Status. Expected: %v. Found: %v",
				"Running", pod.Status.Phase))
		}

		// Validating QDROUTERD_CONF env var in containers
		var envVarsFound []string
		for _, c := range pod.Spec.Containers {
			for _, envVar := range c.Env {
				if len(envVar.Value) == 0 {
					continue
				}
				envVarsFound = append(envVarsFound, envVar.Name)
				switch envVar.Name {
				case "QDROUTERD_CONF":
					if !strings.Contains(envVar.Value, "router {") {
						ginkgo.Fail(fmt.Sprintf("QDROUTERD_CONF does not define the router entity"))
					}
					if !strings.Contains(envVar.Value, "listener {") {
						ginkgo.Fail(fmt.Sprintf("QDROUTERD_CONF does not define any listener"))
					}
				case "APPLICATION_NAME":
					if envVar.Value != cr.ObjectMeta.Name {
						ginkgo.Fail(fmt.Sprintf("APPLICATION_NAME does not match expected value: %v",
							cr.ObjectMeta.Name))
					}
				case "POD_COUNT":
					if envVar.Value != originalReplicaSize {
						ginkgo.Fail(fmt.Sprintf("POD_COUNT does not match expected value: %v. Found: %v",
							originalReplicaSize, envVar.Value))
					}
				}
			}
		}

		if !utils.ContainsAll(utils.FromStrings(expEnvVars), utils.FromStrings(envVarsFound)) {
			ginkgo.Fail(fmt.Sprintf("Missing EnvVars in Pod. Expected: %v. Found: %v",
				expEnvVars, envVarsFound))
		}

	}
	if count != int(cr.Spec.DeploymentPlan.Size) {
		ginkgo.Fail(fmt.Sprintf("Expected pods: %d. Found: %d", int(cr.Spec.DeploymentPlan.Size), count))
	}
}
