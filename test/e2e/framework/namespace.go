// Copyright 2017 The Interconnectedcloud Authors
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
	"sync"
	"time"

	e2elog "github.com/interconnectedcloud/qdr-operator/test/e2e/framework/log"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	// namespaceNamePrefix is the prefix used when creating namespaces with a random name.
	namespaceNamePrefix     = "qdr-operator-e2e-"
	NamespaceCleanupTimeout = 2 * time.Minute
)

func createTestNamespace(client clientset.Interface, name string, labels map[string]string) *corev1.Namespace {
	ginkgo.By(fmt.Sprintf("Creating a namespace %s to execute the test in", name))
	namespace := createNamespace(client, name, labels)
	return namespace
}

func createNamespace(client clientset.Interface, name string, labels map[string]string) *corev1.Namespace {
	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
	}

	namespace, err := client.CoreV1().Namespaces().Create(namespaceObj)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error creating namespace %v", namespaceObj)
	return namespace
}

// CreateNamespace creates a namespace for e2e testing.
func (f *Framework) CreateNamespace(clientSet *clientset.Clientset,
	baseName string, labels map[string]string) *corev1.Namespace {

	ns := createTestNamespace(clientSet, baseName, labels)
	f.AddNamespacesToDelete(ns)
	return ns
}

func (f *Framework) AddNamespacesToDelete(namespaces ...*corev1.Namespace) {
	for _, ns := range namespaces {
		if ns == nil {
			continue
		}
		f.namespacesToDelete = append(f.namespacesToDelete, ns)
	}
}

func generateNamespace(client clientset.Interface, baseName string, labels map[string]string) *corev1.Namespace {
	namespaceObj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("e2e-tests-%v-", baseName),
			Labels:       labels,
		},
	}

	namespace, err := client.CoreV1().Namespaces().Create(namespaceObj)
	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Error generating namespace %v", namespaceObj)
	return namespace
}

// GenerateNamespace creates a namespace with a random name.
func (f *Framework) GenerateNamespace() (*corev1.Namespace, error) {
	return f.KubeClient.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: namespaceNamePrefix,
		},
	})
}

func deleteNamespace(client clientset.Interface, namespaceName string) error {

	return client.CoreV1().Namespaces().Delete(
		namespaceName,
		&metav1.DeleteOptions{})

}

func (f *Framework) DeleteNamespace(ns *corev1.Namespace) []error {
	var errors []error

	if err := deleteNamespace(f.KubeClient, ns.Name); err != nil {
		switch {
		case apierrors.IsNotFound(err):
			e2elog.Logf("Namespace was already deleted")
		case apierrors.IsConflict(err):
			e2elog.Logf("Namespace scheduled for deletion, resources being purged")
		default:
			e2elog.Logf("Failed deleting namespace")
			errors = append(errors, err)
		}
	}

	return errors
}

// DeleteNamespaces deletes all namespaces that match the given delete and skip filters.
// Filter is by simple strings.Contains; first skip filter, then delete filter.
// Returns the list of deleted namespaces or an error.
func DeleteNamespaces(c clientset.Interface, deleteFilter, skipFilter []string) ([]string, error) {
	ginkgo.By("Deleting namespaces")
	nsList, err := c.CoreV1().Namespaces().List(metav1.ListOptions{})
	ExpectNoError(err, "Failed to get namespace list")
	var deleted []string
	var wg sync.WaitGroup
OUTER:
	for _, item := range nsList.Items {
		if skipFilter != nil {
			for _, pattern := range skipFilter {
				if strings.Contains(item.Name, pattern) {
					continue OUTER
				}
			}
		}
		if deleteFilter != nil {
			var shouldDelete bool
			for _, pattern := range deleteFilter {
				if strings.Contains(item.Name, pattern) {
					shouldDelete = true
					break
				}
			}
			if !shouldDelete {
				continue OUTER
			}
		}
		wg.Add(1)
		deleted = append(deleted, item.Name)
		go func(nsName string) {
			defer wg.Done()
			defer ginkgo.GinkgoRecover()
			gomega.Expect(c.CoreV1().Namespaces().Delete(nsName, nil)).To(gomega.Succeed())
			e2elog.Logf("namespace : %v api call to delete is complete ", nsName)
		}(item.Name)
	}
	wg.Wait()
	return deleted, nil
}

// WaitForNamespacesDeleted waits for the namespaces to be deleted.
func WaitForNamespacesDeleted(c clientset.Interface, namespaces []string, timeout time.Duration) error {
	ginkgo.By("Waiting for namespaces to vanish")
	nsMap := map[string]bool{}
	for _, ns := range namespaces {
		nsMap[ns] = true
	}

	return wait.Poll(2*time.Second, timeout,
		func() (bool, error) {
			nsList, err := c.CoreV1().Namespaces().List(metav1.ListOptions{})
			if err != nil {
				return false, err
			}
			for _, item := range nsList.Items {
				if _, ok := nsMap[item.Name]; ok {
					return false, nil
				}
			}
			return true, nil
		})
}
