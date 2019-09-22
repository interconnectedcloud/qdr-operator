package ingresses

import (
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	extv1b1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Create newIngressForCR method to create exposed ingress
func NewIngressForCR(m *v1alpha1.Interconnect, listener v1alpha1.Listener) *extv1b1.Ingress {
	target := listener.Name
	if target == "" {
		target = strconv.Itoa(int(listener.Port))
	}
	labels := selectors.LabelsForInterconnect(m.Name)
	ingress := &extv1b1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      m.Name + "-" + target,
			Namespace: m.Namespace,
		},
		Spec: extv1b1.IngressSpec{
			Rules: []extv1b1.IngressRule{
				{
					Host: m.Name,
					IngressRuleValue: extv1b1.IngressRuleValue{
						HTTP: &extv1b1.HTTPIngressRuleValue{
							Paths: []extv1b1.HTTPIngressPath{
								extv1b1.HTTPIngressPath{
									Path: "/",
									Backend: extv1b1.IngressBackend{
										ServiceName: m.Name,
										ServicePort: intstr.FromInt(int(listener.Port)),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ingress
}
