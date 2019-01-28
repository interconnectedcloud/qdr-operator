package services

import (
	"reflect"
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdrouterd-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/constants"
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/utils/selectors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func nameForListener(l *v1alpha1.Listener) string {
	if l.Name == "" {
		return "port-" + strconv.Itoa(int(l.Port))
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

func portsForQdrouterd(m *v1alpha1.Qdrouterd) []corev1.ServicePort {
	ports := []corev1.ServicePort{}
	external := servicePortsForListeners(m.Spec.Listeners)
	internal := servicePortsForListeners(m.Spec.InterRouterListeners)
	ports = append(ports, external...)
	ports = append(ports, internal...)
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

// Create newServiceForCR method to create service
func NewServiceForCR(m *v1alpha1.Qdrouterd, requestCert bool) *corev1.Service {
	service := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: selectors.LabelsForQdrouterd(m.Name),
			Ports:    portsForQdrouterd(m),
		},
	}
	if requestCert {
		service.Annotations = map[string]string{constants.CertRequestAnnotation: m.Name + "-cert"}
	}
	return service
}
