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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) GetService(name string) (*corev1.Service, error) {
	return f.KubeClient.CoreV1().Services(f.Namespace).Get(name, metav1.GetOptions{})
}

// GetPorts returns an int slice with all ports exposed
// by the provided corev1.Service object
func GetPorts(service corev1.Service) []int {
	if len(service.Spec.Ports) == 0 {
		return []int{}
	}
	var svcPorts []int
	for _, port := range service.Spec.Ports {
		svcPorts = append(svcPorts, int(port.Port))
	}
	return svcPorts
}
