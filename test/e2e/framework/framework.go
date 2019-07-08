// Copyright 2019 The interconnectedcloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package framework

import (
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	qdrclient "github.com/interconnectedcloud/qdr-operator/pkg/client/clientset/versioned"
	e2elog "github.com/interconnectedcloud/qdr-operator/test/e2e/framework/log"
	apiextv1b1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/uuid"
)

const (
	qdrOperatorName = "qdr-operator"
	crdName         = "interconnects.interconnectedcloud.github.io"
	groupName       = "interconnectedcloud.github.io"
	apiVersion      = "v1alpha1"
	opImage         = "quay.io/interconnectedcloud/qdr-operator"
	qdrImage        = "quay.io/interconnectedcloud/qdrouterd:1.8.0"
)

var (
	RetryInterval        = time.Second * 5
	Timeout              = time.Second * 600
	CleanupRetryInterval = time.Second * 1
	CleanupTimeout       = time.Second * 5
	GVR                  = groupName + "/" + apiVersion
)

type Framework struct {
	BaseName string

	// Set together with creating the ClientSet and the namespace.
	// Guaranteed to be unique in the cluster even when running the same
	// test multiple times in parallel.
	UniqueName string

	opImage    string
	KubeClient clientset.Interface
	ExtClient  apiextension.Interface
	QdrClient  qdrclient.Interface

	SkipNamespaceCreation bool // Whether to skip creating a namespace
	Namespace             string
	namespacesToDelete    []*corev1.Namespace // Some tests have more than one
	cleanupHandle         CleanupActionHandle
}

// NewFramework creates a test framework
func NewFramework(baseName string, client clientset.Interface) *Framework {

	f := &Framework{
		BaseName:   baseName,
		KubeClient: client,
	}
	ginkgo.BeforeEach(f.BeforeEach)
	ginkgo.AfterEach(f.AfterEach)

	return f
}

// BeforeEach gets clients and makes a namespace
func (f *Framework) BeforeEach() {

	f.cleanupHandle = AddCleanupAction(f.AfterEach)

	if f.KubeClient == nil {
		ginkgo.By("Creating kubernetes clients")
		config, err := clientcmd.BuildConfigFromFlags("", TestContext.KubeConfig)
		f.KubeClient, err = clientset.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		f.ExtClient, err = apiextension.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		f.QdrClient, err = qdrclient.NewForConfig(config)
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	}
	f.opImage = TestContext.OpImage

	if !f.SkipNamespaceCreation {
		ginkgo.By(fmt.Sprintf("Building namespace api objects, basename %s", f.BaseName))

		namespaceLabels := map[string]string{
			"e2e-framework": f.BaseName,
		}

		namespace := generateNamespace(f.KubeClient, f.BaseName, namespaceLabels)

		f.AddNamespacesToDelete(namespace)
		f.Namespace = namespace.GetName()
		f.UniqueName = namespace.GetName()

	} else {
		f.UniqueName = string(uuid.NewUUID())
	}

	// setup the operator
	err := f.Setup()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}

// AfterEach deletes the namespace, after reading its events.
func (f *Framework) AfterEach() {
	RemoveCleanupAction(f.cleanupHandle)

	// teardown the operator
	err := f.Teardown()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	// DeleteNamespace at the very end in defer, to avoid any
	// expectation failures preventing deleting the namespace.
	defer func() {
		nsDeletionErrors := map[string][]error{}
		// Whether to delete namespace is determined by 3 factors: delete-namespace flag, delete-namespace-on-failure flag and the test result
		// if delete-namespace set to false, namespace will always be preserved.
		// if delete-namespace is true and delete-namespace-on-failure is false, namespace will be preserved if test failed.
		for _, ns := range f.namespacesToDelete {
			ginkgo.By(fmt.Sprintf("Destroying namespace %q for this suite on all clusters.", ns.Name))
			if errors := f.DeleteNamespace(ns); errors != nil {
				nsDeletionErrors[ns.Name] = errors
			}
		}

		// Paranoia-- prevent reuse!
		f.Namespace = ""
		f.KubeClient = nil
		f.namespacesToDelete = nil

		// if we had errors deleting, report them now.
		if len(nsDeletionErrors) != 0 {
			messages := []string{}
			for namespaceKey, namespaceErrors := range nsDeletionErrors {
				for clusterIdx, namespaceErr := range namespaceErrors {
					messages = append(messages, fmt.Sprintf("Couldn't delete ns: %q (@cluster %d): %s (%#v)",
						namespaceKey, clusterIdx, namespaceErr, namespaceErr))
				}
			}
			e2elog.Failf(strings.Join(messages, ","))
		}
	}()

}

func (f *Framework) Teardown() error {

	// Skip the qdr-operator teardown if the operator image was not specified
	if len(f.opImage) == 0 {
		return nil
	}

	err := f.KubeClient.CoreV1().ServiceAccounts(f.Namespace).Delete(qdrOperatorName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator service account: %v", err)
	}
	err = f.KubeClient.RbacV1().Roles(f.Namespace).Delete(qdrOperatorName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator role: %v", err)
	}
	err = f.KubeClient.RbacV1().ClusterRoles().Delete(qdrOperatorName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator cluster role: %v", err)
	}
	err = f.KubeClient.RbacV1().RoleBindings(f.Namespace).Delete(qdrOperatorName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator role binding: %v", err)
	}
	err = f.KubeClient.RbacV1().ClusterRoleBindings().Delete(qdrOperatorName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator cluster role binding: %v", err)
	}
	err = f.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crdName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator crd: %v", err)
	}
	err = f.KubeClient.AppsV1().Deployments(f.Namespace).Delete(qdrOperatorName, metav1.NewDeleteOptions(1))
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to delete qdr-operator deployment: %v", err)
	}

	e2elog.Logf("e2e teardown succesful")
	return nil
}

func (f *Framework) Setup() error {

	err := f.setupQdrServiceAccount()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = f.setupQdrRole()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = f.setupQdrClusterRole()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = f.setupQdrRoleBinding()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = f.setupQdrClusterRoleBinding()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = f.setupQdrCrd()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = f.setupQdrDeployment()
	if err != nil {
		return fmt.Errorf("failed to setup qdr operator: %v", err)
	}
	err = WaitForDeployment(f.KubeClient, f.Namespace, "qdr-operator", 1, RetryInterval, Timeout)
	if err != nil {
		return fmt.Errorf("Failed to wait for qdr operator: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrServiceAccount() error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: qdrOperatorName,
		},
	}
	_, err := f.KubeClient.CoreV1().ServiceAccounts(f.Namespace).Create(sa)
	if err != nil {
		return fmt.Errorf("create qdr-operator service account failed: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrRole() error {
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name: qdrOperatorName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services", "serviceaccounts", "endpoints", "persistentvolumeclaims", "events", "configmaps", "secrets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"rolebindings", "roles"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{"extensions"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"namespaces"},
				Verbs:     []string{"get"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments", "daemonsets", "replicasets", "statefulsets"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{"certmanager.k8s.io"},
				Resources: []string{"issuers", "certificates"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{"monitoring.coreos.com"},
				Resources: []string{"servicemonitors"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{"route.openshift.io"},
				Resources: []string{"routes", "routes/custom-host", "routes/status"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{"interconnectedcloud.github.io"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
		},
	}
	_, err := f.KubeClient.RbacV1().Roles(f.Namespace).Create(role)
	if err != nil {
		return fmt.Errorf("create qdr-operator role failed: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrClusterRole() error {
	crole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: qdrOperatorName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"apiextensions.k8s.io"},
				Resources: []string{"customresourcedefinitions"},
				Verbs:     []string{"get", "list"},
			},
		},
	}
	_, err := f.KubeClient.RbacV1().ClusterRoles().Create(crole)
	if err != nil {
		return fmt.Errorf("create qdr-operator cluster role failed: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrRoleBinding() error {
	rb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: qdrOperatorName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     qdrOperatorName,
		},
		Subjects: []rbacv1.Subject{
			{
				APIGroup:  "",
				Kind:      "ServiceAccount",
				Name:      qdrOperatorName,
				Namespace: f.Namespace,
			},
		},
	}
	_, err := f.KubeClient.RbacV1().RoleBindings(f.Namespace).Create(rb)
	if err != nil {
		return fmt.Errorf("create qdr-operator role binding failed: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrClusterRoleBinding() error {
	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: qdrOperatorName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     qdrOperatorName,
		},
		Subjects: []rbacv1.Subject{
			{
				APIGroup:  "",
				Kind:      "ServiceAccount",
				Name:      qdrOperatorName,
				Namespace: f.Namespace,
			},
		},
	}
	_, err := f.KubeClient.RbacV1().ClusterRoleBindings().Create(crb)
	if err != nil {
		return fmt.Errorf("create qdr-operator cluster role binding failed: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrCrd() error {
	crd := &apiextv1b1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdName,
		},
		Spec: apiextv1b1.CustomResourceDefinitionSpec{
			Group: "interconnectedcloud.github.io",
			Names: apiextv1b1.CustomResourceDefinitionNames{
				Kind:     "Interconnect",
				ListKind: "InterconnectList",
				Plural:   "interconnects",
				Singular: "interconnect",
			},
			Scope:   "Namespaced",
			Version: "v1alpha1",
		},
	}
	_, err := f.ExtClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		return fmt.Errorf("create qdr-operator crd failed: %v", err)
	}
	return nil
}

func (f *Framework) setupQdrDeployment() error {
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: qdrOperatorName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": qdrOperatorName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": qdrOperatorName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: qdrOperatorName,
					Containers: []corev1.Container{
						{
							Command:         []string{qdrOperatorName},
							Name:            qdrOperatorName,
							Image:           f.opImage,
							ImagePullPolicy: corev1.PullAlways,
							Env: []corev1.EnvVar{
								{
									Name:      "WATCH_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"}},
								},
								{
									Name:      "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.name"}},
								},
								{
									Name:  "OPERATOR_NAME",
									Value: qdrOperatorName,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 60000,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := f.KubeClient.AppsV1().Deployments(f.Namespace).Create(dep)
	if err != nil {
		return fmt.Errorf("create qdr-operator deployment failed: %v", err)
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
