/*
Copyright 2019 The Interconnectedcloud Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"fmt"
	"sort"

	"k8s.io/apimachinery/pkg/util/uuid"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	e2elog "github.com/interconnectedcloud/qdr-operator/test/e2e/framework/log"

	"github.com/onsi/gomega"
)

var (
	RunID = uuid.NewUUID()
)

// RestclientConfig returns a config holds the information needed to build connection to kubernetes clusters.
func RestclientConfig(kubeContext string) (*clientcmdapi.Config, error) {
	//e2elog.Logf(">>> kubeConfig: %s", TestContext.KubeConfig)
	if TestContext.KubeConfig == "" {
		return nil, fmt.Errorf("KubeConfig must be specified to load client config")
	}
	c, err := clientcmd.LoadFromFile(TestContext.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("error loading KubeConfig: %v", err.Error())
	}
	if kubeContext != "" {
		//e2elog.Logf(">>> kubeContext: %s", kubeContext)
		c.CurrentContext = kubeContext
	}
	return c, nil
}

// LoadConfig returns a config for a rest client.
func LoadConfig() (*restclient.Config, error) {
	c, err := RestclientConfig(TestContext.KubeContext)
	if err != nil {
		if TestContext.KubeConfig == "" {
			return restclient.InClusterConfig()
		}
		return nil, err
	}

	return clientcmd.NewDefaultClientConfig(*c, &clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: TestContext.Host}}).ClientConfig()
}

// LoadClientset returns clientset for connecting to kubernetes clusters.
func LoadClientset() (*clientset.Clientset, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err.Error())
	}
	return clientset.NewForConfig(config)
}

// ExpectError expects an error happens, otherwise an exception raises
func ExpectError(err error, explain ...interface{}) {
	gomega.Expect(err).To(gomega.HaveOccurred(), explain...)
}

// ExpectNoError checks if "err" is set, and if so, fails assertion while logging the error.
func ExpectNoError(err error, explain ...interface{}) {
	ExpectNoErrorWithOffset(1, err, explain...)
}

// ExpectNoErrorWithOffset checks if "err" is set, and if so, fails assertion while logging the error at "offset" levels above its caller
// (for example, for call chain f -> g -> ExpectNoErrorWithOffset(1, ...) error would be logged for "f").
func ExpectNoErrorWithOffset(offset int, err error, explain ...interface{}) {
	if err != nil {
		e2elog.Logf("Unexpected error occurred: %v", err)
	}
	gomega.ExpectWithOffset(1+offset, err).NotTo(gomega.HaveOccurred(), explain...)
}

// ExpectNoErrorWithRetries checks if an error occurs with the given retry count.
func ExpectNoErrorWithRetries(fn func() error, maxRetries int, explain ...interface{}) {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = fn()
		if err == nil {
			return
		}
		e2elog.Logf("(Attempt %d of %d) Unexpected error occurred: %v", i+1, maxRetries, err)
	}
	gomega.ExpectWithOffset(1, err).NotTo(gomega.HaveOccurred(), explain...)
}

// ContainsAll can be used to compare two sorted slices and validate
// both are not nil and all elements from given model are present on
// the target instance.
func ContainsAll(model, target []interface{}) bool {
	if model == nil || target == nil {
		return false
	}

	if len(model) > len(target) {
		return false
	}

	ti := 0
	for _, vm := range model {
		found := false
		for i := ti; i < len(target); i++ {
			if vm == target[i] {
				found = true
				if ti == i {
					ti++
				} else {
					ti = i
				}
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// FromInts returns an interface array with sorted elements from int array
func FromInts(m []int) []interface{} {
	res := make([]interface{}, len(m))
	for i, v := range m {
		res[i] = v
	}
	sort.Slice(res, func(i1, i2 int) bool {
		return res[i1].(int) < res[i2].(int)
	})
	return res
}

// FromStrings returns an interface array with sorted elements from string array
func FromStrings(m []string) []interface{} {
	res := make([]interface{}, len(m))
	for i, v := range m {
		res[i] = v
	}
	sort.Slice(res, func(i1, i2 int) bool {
		return res[i1].(string) < res[i2].(string)
	})
	return res
}
