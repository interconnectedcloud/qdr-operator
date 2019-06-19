package validations

import (
	"context"
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/utils"
	"github.com/onsi/ginkgo"
	"k8s.io/api/rbac/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ValidateRoles will retrieve the RoleList available in the
// current namespace defined in the CR and validate that
// a Role with the CR Name exists and another Role for
// the qdr-operator.
func ValidateRoles(cr *v1alpha1.Interconnect) {

	// Roles that must exist in the namespace
	expectedRoles := []string{cr.ObjectMeta.Name, "qdr-operator"}

	cliListOptions := client.ListOptions{
		Namespace: cr.ObjectMeta.Namespace,
	}
	roleList := v1.RoleList{
		TypeMeta: v12.TypeMeta{
			Kind:       "Role",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
	}
	err := e2e.TestSuite.F.Client.List(context.TODO(), &cliListOptions, &roleList)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	// Retrieving all roles found
	var rolesFound []string
	if len(roleList.Items) > 0 {
		for _, role := range roleList.Items {
			rolesFound = append(rolesFound, role.Name)
		}
	}

	// If roles found do not match expected
	if !utils.ContainsAll(utils.FromStrings(expectedRoles), utils.FromStrings(rolesFound)) {
		ginkgo.Fail(fmt.Sprintf("Expected: %v. Found: %v.", expectedRoles, rolesFound))
	}

}

// ValidateRoleBindings verifies that namespace contains at least two
// RoleBinding instances defined. One named as qdr-operator and the
// other must have the CR name.
func ValidateRoleBindings(cr *v1alpha1.Interconnect) {
	cliListOptions := client.ListOptions{
		Namespace: cr.ObjectMeta.Namespace,
	}
	roleBindingList := v1.RoleBindingList{
		TypeMeta: v12.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
	}
	err := e2e.TestSuite.F.Client.List(context.TODO(), &cliListOptions, &roleBindingList)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
	// Expects that two role bindings are defined
	// One for the given CR name and a static qdr-operator
	expectedRb := []string{
		cr.ObjectMeta.Name,
		"qdr-operator",
	}
	rbFound := []string{}
	if len(roleBindingList.Items) > 0 {
		for _, rb := range roleBindingList.Items {
			rbFound = append(rbFound, rb.Name)
		}
	}
	if !utils.ContainsAll(utils.FromStrings(expectedRb), utils.FromStrings(rbFound)) {
		ginkgo.Fail(fmt.Sprintf("Expected: %v. Found: %v", expectedRb, rbFound))
	}
}
