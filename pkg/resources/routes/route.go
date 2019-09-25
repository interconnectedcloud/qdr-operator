package routes

import (
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Create newRouteForCR method to create exposed route
func NewRouteForCR(m *v1alpha1.Interconnect, listener v1alpha1.Listener) *routev1.Route {
	target := listener.Name
	if target == "" {
		target = strconv.Itoa(int(listener.Port))
	}
	labels := selectors.LabelsForInterconnect(m.Name)
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
				TargetPort: intstr.FromInt(int(listener.Port)),
			},
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: m.Name,
			},
		},
	}
	if listener.SslProfile != "" {
		route.Spec.TLS = &routev1.TLSConfig{
			Termination:                   routev1.TLSTerminationPassthrough,
			InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyNone,
		}
	} else {
		route.Spec.TLS = &routev1.TLSConfig{
			Termination:                   routev1.TLSTerminationEdge,
			InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
		}
	}
	return route
}
