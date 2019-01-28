package containers

import (
	"reflect"
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdrouterd-operator/pkg/apis/interconnectedcloud/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func containerPortsForListeners(listeners []v1alpha1.Listener) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{}
	for _, listener := range listeners {
		ports = append(ports, corev1.ContainerPort{
			Name:          nameForListener(&listener),
			ContainerPort: listener.Port,
		})
	}
	return ports
}

func containerPortsForQdrouterd(m *v1alpha1.Qdrouterd) []corev1.ContainerPort {
	ports := containerPortsForListeners(m.Spec.Listeners)
	ports = append(ports, containerPortsForListeners(m.Spec.InterRouterListeners)...)
	return ports
}

func nameForListener(l *v1alpha1.Listener) string {
	if l.Name == "" {
		return "port-" + strconv.Itoa(int(l.Port))
	} else {
		return l.Name
	}
}

func CheckQdrouterdContainer(desired *corev1.Container, actual *corev1.Container) bool {
	if desired.Image != actual.Image {
		return false
	}
	if !reflect.DeepEqual(desired.Env, actual.Env) {
		return false
	}
	if !reflect.DeepEqual(desired.Ports, actual.Ports) {
		return false
	}
	if !reflect.DeepEqual(desired.VolumeMounts, actual.VolumeMounts) {
		return false
	}
	return true
}

func ContainerForQdrouterd(m *v1alpha1.Qdrouterd, config string) corev1.Container {
	container := corev1.Container{
		Image: m.Spec.Image,
		Name:  m.Name,
		Env: []corev1.EnvVar{
			{
				Name:  "QDROUTERD_CONF",
				Value: config,
			},
			{
				Name:  "QDROUTERD_AUTO_MESH_DISCOVERY",
				Value: "QUERY",
			},
			{
				Name:  "APPLICATION_NAME",
				Value: m.Name,
			},
			{
				Name: "POD_NAMESPACE",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "metadata.namespace",
					},
				},
			},
			{
				Name: "POD_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
		},
		Ports: containerPortsForQdrouterd(m),
	}
	if m.Spec.SslProfiles != nil && len(m.Spec.SslProfiles) > 0 {
		volumeMounts := []corev1.VolumeMount{}
		for _, profile := range m.Spec.SslProfiles {
			if len(profile.Credentials) > 0 {
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      profile.Credentials,
					MountPath: "/etc/qpid-dispatch-certs/" + profile.Name + "/" + profile.Credentials,
				})
			}
			if len(profile.CaCert) > 0 && profile.CaCert != profile.Credentials {
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      profile.CaCert,
					MountPath: "/etc/qpid-dispatch-certs/" + profile.Name + "/" + profile.CaCert,
				})
			}

		}
		container.VolumeMounts = volumeMounts
	}
	return container
}
