package routes

import (
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Create newRouteForCR method to create exposed route
func NewRouteForCR(m *v1alpha1.Qdr, target string) *routev1.Route {
	labels := selectors.LabelsForQdr(m.Name)
	route := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Route",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      m.Name + "-" + target,
			Namespace: m.Namespace,
		},
		Spec: routev1.RouteSpec{
			Path: "",
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString(target),
			},
			TLS: &routev1.TLSConfig{
				Termination:                   routev1.TLSTerminationPassthrough,
				InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
			},
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: m.Name,
			},
		},
	}
	return route
}
