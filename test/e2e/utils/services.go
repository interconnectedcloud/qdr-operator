package utils

import "k8s.io/api/core/v1"

// GetPorts returns an int slice with all ports exposed
// by the provided corev1.Service object
func GetPorts(service v1.Service) []int {
	if len(service.Spec.Ports) == 0 {
		return []int{}
	}
	var svcPorts []int
	for _, port := range service.Spec.Ports {
		svcPorts = append(svcPorts, int(port.Port))
	}
	return svcPorts
}
