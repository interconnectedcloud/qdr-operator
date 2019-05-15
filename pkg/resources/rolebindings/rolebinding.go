package rolebindings

import (
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Create NewRoleBindingForCR method to create rolebinding
func NewRoleBindingForCR(m *v1alpha1.Interconnect) *rbacv1.RoleBinding {
	rolebinding := &rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "rbac.authorization.k8s.io/v1",
			Kind:       "RoleBinding",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Subjects: []rbacv1.Subject{{
			Kind: "ServiceAccount",
			Name: m.Name,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind: "Role",
			Name: m.Name,
		},
	}

	return rolebinding
}
