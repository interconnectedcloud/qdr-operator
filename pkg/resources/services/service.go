package services

import (
	"reflect"
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/constants"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	log = logf.Log.WithName("services")
)

func nameForListener(l *v1alpha1.Listener) string {
	if l.Name == "" {
		return strconv.Itoa(int(l.Port))
	} else {
		return l.Name
	}
}

func servicePortsForListeners(listeners []v1alpha1.Listener) []corev1.ServicePort {
	ports := []corev1.ServicePort{}
	for _, listener := range listeners {
		ports = append(ports, corev1.ServicePort{
			Name:       nameForListener(&listener),
			Protocol:   "TCP",
			Port:       listener.Port,
			TargetPort: intstr.FromInt(int(listener.Port)),
		})
	}
	return ports
}

func servicePortsForRouter(m *v1alpha1.Interconnect) []corev1.ServicePort {
	ports := []corev1.ServicePort{}
	external := servicePortsForListeners(m.Spec.Listeners)
	internal := servicePortsForListeners(m.Spec.InterRouterListeners)
	edge := servicePortsForListeners(m.Spec.EdgeListeners)
	ports = append(ports, external...)
	ports = append(ports, internal...)
	ports = append(ports, edge...)
	return ports
}

func CheckService(desired *corev1.Service, actual *corev1.Service) bool {
	update := false
	if !reflect.DeepEqual(desired.Annotations[constants.CertRequestAnnotation], actual.Annotations[constants.CertRequestAnnotation]) {
		actual.Annotations[constants.CertRequestAnnotation] = desired.Annotations[constants.CertRequestAnnotation]
	}
	if !reflect.DeepEqual(desired.Spec.Selector, actual.Spec.Selector) {
		actual.Spec.Selector = desired.Spec.Selector
	}
	if !reflect.DeepEqual(desired.Spec.Ports, actual.Spec.Ports) {
		actual.Spec.Ports = desired.Spec.Ports
	}
	return update
}

// Create newServiceForCR method to create normal service
func NewServiceForCR(m *v1alpha1.Interconnect, requestCert bool) *corev1.Service {
	labels := selectors.LabelsForInterconnect(m.Name)
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    servicePortsForRouter(m),
		},
	}
	if requestCert {
		service.Annotations = map[string]string{constants.CertRequestAnnotation: m.Name + "-cert"}
	}
	if m.Spec.DeploymentPlan.ServiceType == "ClusterIP" || m.Spec.DeploymentPlan.ServiceType == "NodePort" || m.Spec.DeploymentPlan.ServiceType == "LoadBalancer" {
		service.Spec.Type = corev1.ServiceType(m.Spec.DeploymentPlan.ServiceType)
	} else if m.Spec.DeploymentPlan.ServiceType != "" {
		log.Info("Unrecognised ServiceType", "ServiceType", m.Spec.DeploymentPlan.ServiceType)
	}
	return service
}

// Create newServiceForCR method to create normal service
func NewNormalServiceForCR(m *v1alpha1.Interconnect, requestCert bool) *corev1.Service {
	labels := selectors.LabelsForInterconnect(m.Name)
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      m.Name + "-normal",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Type:     "LoadBalancer",
			Selector: labels,
			Ports:    servicePortsForListeners(m.Spec.Listeners),
		},
	}
	if requestCert {
		service.Annotations = map[string]string{constants.CertRequestAnnotation: m.Name + "-cert"}
	}
	return service
}

// Create newHeadlessServiceForCR method to create normal service
func NewHeadlessServiceForCR(m *v1alpha1.Interconnect, requestCert bool) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name + "-headless",
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Selector:  selectors.LabelsForInterconnect(m.Name),
			Ports:     servicePortsForListeners(m.Spec.InterRouterListeners),
		},
	}
	if requestCert {
		service.Annotations = map[string]string{constants.CertRequestAnnotation: m.Name + "-cert"}
	}
	return service
}
