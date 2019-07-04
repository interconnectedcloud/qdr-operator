package deployments

import (
	"strconv"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/resources/containers"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// move this to util
// Set labels in a map
func labelsForInterconnect(name string) map[string]string {
	return map[string]string{
		selectors.LabelAppKey:      name,
		selectors.LabelResourceKey: name,
	}
}

func CheckDeployedContainer(actual *corev1.PodTemplateSpec, cr *v1alpha1.Interconnect) bool {
	desired := containers.ContainerForInterconnect(cr)
	if len(actual.Spec.Containers) != 1 || !containers.CheckInterconnectContainer(&desired, &actual.Spec.Containers[0]) {
		actual.Spec.Containers = []corev1.Container{desired}
		return false
	}
	return true
}

// Create NewDeploymentForCR method to create deployment
func NewDeploymentForCR(m *v1alpha1.Interconnect) *appsv1.Deployment {
	labels := selectors.LabelsForInterconnect(m.Name)
	replicas := m.Spec.DeploymentPlan.Size
	affinity := &corev1.Affinity{}
	if m.Spec.DeploymentPlan.Placement == v1alpha1.PlacementAntiAffinity {
		affinity = &corev1.Affinity{
			PodAntiAffinity: &corev1.PodAntiAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
					{
						LabelSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "application",
									Operator: metav1.LabelSelectorOpIn,
									Values:   []string{m.Name},
								},
							},
						},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		}
	}
	dep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						"prometheus.io/port":   strconv.Itoa(int(m.Spec.DeploymentPlan.LivenessPort)),
						"prometheus.io/scrape": "true",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: m.Name,
					Affinity:           affinity,
					Containers:         []corev1.Container{containers.ContainerForInterconnect(m)},
				},
			},
		},
	}
	volumes := []corev1.Volume{}
	for _, profile := range m.Spec.SslProfiles {
		if len(profile.Credentials) > 0 {
			volumes = append(volumes, corev1.Volume{
				Name: profile.Credentials,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: profile.Credentials,
					},
				},
			})
		}
		if len(profile.CaCert) > 0 && profile.CaCert != profile.Credentials {
			volumes = append(volumes, corev1.Volume{
				Name: profile.CaCert,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: profile.CaCert,
					},
				},
			})
		}
	}
	if len(m.Spec.Users) > 0 {
		volumes = append(volumes, corev1.Volume{
			Name: "sasl-users",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: m.Spec.Users,
				},
			},
		})
		volumes = append(volumes, corev1.Volume{
			Name: "sasl-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.Name + "-sasl-config",
					},
				},
			},
		})
	}

	dep.Spec.Template.Spec.Volumes = volumes

	return dep
}

// Create NewDaemonSetForCR method to create daemonset
func NewDaemonSetForCR(m *v1alpha1.Interconnect) *appsv1.DaemonSet {
	labels := selectors.LabelsForInterconnect(m.Name)

	ds := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Annotations: map[string]string{
						"prometheus.io/port":   strconv.Itoa(int(m.Spec.DeploymentPlan.LivenessPort)),
						"prometheus.io/scrape": "true",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: m.Name,
					Containers:         []corev1.Container{containers.ContainerForInterconnect(m)},
				},
			},
		},
	}
	volumes := []corev1.Volume{}
	for _, profile := range m.Spec.SslProfiles {
		if len(profile.Credentials) > 0 {
			volumes = append(volumes, corev1.Volume{
				Name: profile.Credentials,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: profile.Credentials,
					},
				},
			})
		}
		if len(profile.CaCert) > 0 && configs.IsCaSecretNeeded(&profile) {
			volumes = append(volumes, corev1.Volume{
				Name: profile.CaCert,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: profile.CaCert,
					},
				},
			})
		}
	}
	if len(m.Spec.Users) > 0 {
		volumes = append(volumes, corev1.Volume{
			Name: "sasl-users",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: m.Spec.Users,
				},
			},
		})
		volumes = append(volumes, corev1.Volume{
			Name: "sasl-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.Name + "-sasl-config",
					},
				},
			},
		})
	}

	ds.Spec.Template.Spec.Volumes = volumes

	return ds
}
