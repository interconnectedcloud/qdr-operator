package containers

import (
	"os"
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	log = logf.Log.WithName("containers")
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

func containerPortsForInterconnect(m *v1alpha1.Interconnect) []corev1.ContainerPort {
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

func containerEnvVarsForInterconnect(m *v1alpha1.Interconnect) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}
	envVars = append(envVars, corev1.EnvVar{Name: "APPLICATION_NAME", Value: m.Name})
	envVars = append(envVars, corev1.EnvVar{Name: "QDROUTERD_CONF", Value: configs.ConfigForInterconnect(m)})
	if m.Spec.Users != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "QDROUTERD_AUTO_CREATE_SASLDB_SOURCE", Value: "/etc/qpid-dispatch/sasl-users/"})
		envVars = append(envVars, corev1.EnvVar{Name: "QDROUTERD_AUTO_CREATE_SASLDB_PATH", Value: "/tmp/qdrouterd.sasldb"})
	}
	envVars = append(envVars, corev1.EnvVar{Name: "POD_COUNT", Value: strconv.Itoa(int(m.Spec.DeploymentPlan.Size))})
	envVars = append(envVars, corev1.EnvVar{Name: "POD_NAMESPACE", ValueFrom: &corev1.EnvVarSource{
		FieldRef: &corev1.ObjectFieldSelector{
			FieldPath: "metadata.namespace",
		},
	},
	})
	envVars = append(envVars, corev1.EnvVar{Name: "POD_IP", ValueFrom: &corev1.EnvVarSource{
		FieldRef: &corev1.ObjectFieldSelector{
			FieldPath: "status.podIP",
		},
	},
	})
	if m.Spec.DeploymentPlan.Role == v1alpha1.RouterRoleInterior {
		envVars = append(envVars, corev1.EnvVar{Name: "QDROUTERD_AUTO_MESH_DISCOVERY", Value: "QUERY"})
	}
	return envVars
}

func findEnvVar(env []corev1.EnvVar, name string) *corev1.EnvVar {
	for _, v := range env {
		if v.Name == name {
			return &v
		}
	}
	return nil
}

func checkInterconnectConfig(desired []corev1.EnvVar, actual []corev1.EnvVar) bool {
	a := findEnvVar(desired, "QDROUTERD_CONF")
	b := findEnvVar(actual, "QDROUTERD_CONF")
	return (a == nil && b == nil) || a.Value == b.Value
}

func CheckInterconnectContainer(desired *corev1.Container, actual *corev1.Container) bool {
	if desired.Image != actual.Image {
		log.Info("Name changed", "desired", desired.Image, "actual", actual.Image)
		return false
	}
	if !checkInterconnectConfig(desired.Env, actual.Env) {
		log.Info("Config changed", "desired", desired.Env, "actual", actual.Env)
		return false
	}
	return true
}

func ContainerForInterconnect(m *v1alpha1.Interconnect) corev1.Container {
	var image string
	if m.Spec.DeploymentPlan.Image != "" {
		image = m.Spec.DeploymentPlan.Image
	} else {
		image = os.Getenv("QDROUTERD_IMAGE")
	}
	container := corev1.Container{
		Image: image,
		Name:  m.Name,
		LivenessProbe: &corev1.Probe{
			InitialDelaySeconds: 60,
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Port: intstr.FromInt(int(m.Spec.DeploymentPlan.LivenessPort)),
					Path: "/healthz",
				},
			},
		},
		Env:   containerEnvVarsForInterconnect(m),
		Ports: containerPortsForInterconnect(m),
	}
	volumeMounts := []corev1.VolumeMount{}
	if m.Spec.SslProfiles != nil && len(m.Spec.SslProfiles) > 0 {
		for _, profile := range m.Spec.SslProfiles {
			if len(profile.Credentials) > 0 {
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      profile.Credentials,
					MountPath: "/etc/qpid-dispatch-certs/" + profile.Name + "/" + profile.Credentials,
				})
			}
			if len(profile.CaCert) > 0 && configs.IsCaSecretNeeded(&profile) {
				volumeMounts = append(volumeMounts, corev1.VolumeMount{
					Name:      profile.CaCert,
					MountPath: "/etc/qpid-dispatch-certs/" + profile.Name + "/" + profile.CaCert,
				})
			}

		}
	}
	if len(m.Spec.Users) > 0 {
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "sasl-users",
			MountPath: "/etc/qpid-dispatch/sasl-users",
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      "sasl-config",
			MountPath: "/etc/sasl2",
		})
	}

	container.VolumeMounts = volumeMounts
	container.Resources = m.Spec.DeploymentPlan.Resources
	return container
}
