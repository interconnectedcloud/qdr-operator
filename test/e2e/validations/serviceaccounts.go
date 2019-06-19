package validations

import (
	"context"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/utils"
	"github.com/onsi/ginkgo"
	"k8s.io/api/core/v1"
	v13 "k8s.io/api/rbac/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ValidateServiceAccounts verifies that namespace contains at least two
// ServiceAccount instances defined. One named as qdr-operator and the
// other must have the CR name.
func ValidateServiceAccounts(cr *v1alpha1.Interconnect) {
	cliListOptions := client.ListOptions{
		Namespace: cr.ObjectMeta.Namespace,
	}
	svcAccountList := v1.ServiceAccountList{
		TypeMeta: v12.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: v13.SchemeGroupVersion.String(),
		},
	}
	err := e2e.TestSuite.F.Client.List(context.TODO(), &cliListOptions, &svcAccountList)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	expectedSvcAccounts := []string{
		cr.ObjectMeta.Name,
		"qdr-operator",
	}
	var svcAccountsFound []string
	if len(svcAccountList.Items) > 0 {
		for _, svcAccount := range svcAccountList.Items {
			svcAccountsFound = append(svcAccountsFound, svcAccount.Name)
		}
	}
	if !utils.ContainsAll(utils.FromStrings(expectedSvcAccounts), utils.FromStrings(svcAccountsFound)) {
		ginkgo.Fail(fmt.Sprint("Expected", expectedSvcAccounts, "Found", svcAccountsFound))
	}
}
