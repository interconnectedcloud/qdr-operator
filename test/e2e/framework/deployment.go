// Copyright 2018 The Operator-SDK Authors
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
	"context"
	"testing"
	"time"

	e2elog "github.com/interconnectedcloud/qdr-operator/test/e2e/framework/log"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (f *Framework) GetDeployment(name string) (*appsv1.Deployment, error) {
	return f.KubeClient.AppsV1().Deployments(f.Namespace).Get(name, metav1.GetOptions{})
}

func WaitForDeployment(kubeclient kubernetes.Interface, namespace, name string, replicas int, retryInterval, timeout time.Duration) error {
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		deployment, err := kubeclient.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{IncludeUninitialized: true})
		if err != nil {
			if apierrors.IsNotFound(err) {
				e2elog.Logf("Waiting for availability of %s deployment\n", name)
				return false, nil
			}
			return false, err
		}

		if int(deployment.Status.AvailableReplicas) == replicas {
			return true, nil
		}
		e2elog.Logf("Waiting for full availability of %s deployment (%d/%d)\n", name, deployment.Status.AvailableReplicas, replicas)
		return false, nil
	})
	if err != nil {
		return err
	}
	e2elog.Logf("Deployment available (%d/%d)\n", replicas, replicas)
	return nil
}

func WaitForDeletion(t *testing.T, dynclient client.Client, obj runtime.Object, retryInterval, timeout time.Duration) error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return err
	}

	kind := obj.GetObjectKind().GroupVersionKind().Kind
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		err = dynclient.Get(ctx, key, obj)
		if apierrors.IsNotFound(err) {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		e2elog.Logf("Waiting for %s %s to be deleted\n", kind, key)
		return false, nil
	})
	if err != nil {
		return err
	}
	e2elog.Logf("%s %s was deleted\n", kind, key)
	return nil
}
