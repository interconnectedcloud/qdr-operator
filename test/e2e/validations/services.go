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
)

// ValidateService ensures that there is a Service object in the current
// namespace, it is owned by the provided CR and expectedPorts match
// the ports available through the found Service object.
func ValidateService(cr *v1alpha1.Interconnect, expectedPorts []int) {
	cliListOptions := client.ListOptions{
		Namespace: cr.ObjectMeta.Namespace,
	}
	svcList := v1.ServiceList{
		TypeMeta: v12.TypeMeta{
			Kind:       "Pod",
			APIVersion: v13.SchemeGroupVersion.String(),
		},
	}
	err := e2e.TestSuite.F.Client.List(context.TODO(), &cliListOptions, &svcList)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	if len(svcList.Items) == 0 {
		ginkgo.Fail(fmt.Sprintf("No Services found"))
	}
	found := false
	for _, svc := range svcList.Items {

		if svc.Name != cr.ObjectMeta.Name {
			continue
		}
		found = true

		// Validate owner for given service is the CR
		if cr.ObjectMeta.Name != svc.ObjectMeta.OwnerReferences[0].Name {
			ginkgo.Fail(fmt.Sprintf("Invalid service owner. Expected: %v. Got: %v",
				cr.ObjectMeta.Name,
				svc.ObjectMeta.OwnerReferences[0].Name))
		}

		// Validate expected ports
		portsFound := utils.GetPorts(svc)
		if !utils.ContainsAll(utils.FromInts(expectedPorts), utils.FromInts(portsFound)) {
			ginkgo.Fail(fmt.Sprintf("Expected ports not available. Expected: %v. Found: %v",
				expectedPorts,
				portsFound))
		}

		break
	}
	if !found {
		ginkgo.Fail(fmt.Sprintf("Service not found. Expected: %v.", cr.ObjectMeta.Name))
	}
}
