package components

import (
	"github.com/interconnectedcloud/qdr-operator/tools/constants"
	"strings"

	monv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Verbs = []string{
	"create",
	"delete",
	"deletecollection",
	"get",
	"list",
	"patch",
	"update",
	"watch",
}

func GetDeployment(operatorName, repository, context, imageName, tag, imagePullPolicy string) *appsv1.Deployment {
	registryName := strings.Join([]string{repository, context, imageName}, "/")
	image := strings.Join([]string{registryName, tag}, ":")
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			/*Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxSurge:       &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
					MaxUnavailable: &intstr.IntOrString{Type: intstr.Int, IntVal: 1},
				},
			},*/
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": operatorName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": operatorName,
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: operatorName,
					Containers: []corev1.Container{
						{
							Name:  operatorName,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "metrics",
									ContainerPort: 60000,
								},
							},
							ImagePullPolicy: corev1.PullPolicy(imagePullPolicy),
							Command:         []string{operatorName},
							Env: []corev1.EnvVar{
								{
									Name: "WATCH_NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
								{
									Name: "POD_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
								{
									Name:  "OPERATOR_NAME",
									Value: operatorName,
								},
							},
						},
					},
				},
			},
		},
	}

	deployment.Spec.Template.Spec.Containers[0].Env = append(deployment.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
		Name:  constants.InterconnectVar,
		Value: constants.Interconnect17ImageURL,
	})

	return deployment
}

func GetRole(operatorName string) *rbacv1.Role {
	role := &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			APIVersion: rbacv1.SchemeGroupVersion.String(),
			Kind:       "Role",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: operatorName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"pods",
					"services",
					"serviceaccounts",
					"endpoints",
					"persistentvolumeclaims",
					"events",
					"configmaps",
					"secrets",
				},
				Verbs: []string{"*"},
			},
			{
				APIGroups: []string{
					"rbac.authorization.k8s.io",
				},
				Resources: []string{
					"rolebindings",
					"roles",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"delete"},
			},
			{
				APIGroups: []string{
					"extensions",
				},
				Resources: []string{
					"ingresses",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"delete"},
			},
			{
				APIGroups: []string{
					"",
				},
				Resources: []string{
					"namespaces",
				},
				Verbs: []string{"get"},
			},
			{
				APIGroups: []string{
					"apps",
				},
				Resources: []string{
					"deployments",
					"daemonsets",
					"replicasets",
					"statefulsets",
				},
				Verbs: []string{"*"},
			},

			{
				APIGroups: []string{
					"certmanager.k8s.io",
				},
				Resources: []string{
					"issuers",
					"certificates",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"delete"},
			},
			{
				APIGroups: []string{
					monv1.SchemeGroupVersion.Group,
				},
				Resources: []string{"servicemonitors"},
				Verbs:     []string{"get", "create"},
			},
			{
				APIGroups: []string{
					"route.openshift.io",
				},
				Resources: []string{
					"routes",
					"routes/custom-host",
					"routes/status",
				},
				Verbs: []string{
					"get",
					"list",
					"watch",
					"create",
					"delete"},
			},
			{
				APIGroups: []string{
					"interconnectedcloud.github.io",
				},
				Resources: []string{
					"*",
				},
				Verbs: []string{"*"},
			},
		},
	}
	return role
}

func int32Ptr(i int32) *int32 {
	return &i
}

//GetImage
func GetImage(imageURL string) (image, imageTag, imageContext string) {
	urlParts := strings.Split(imageURL, "/")
	if len(urlParts) > 1 {
		imageContext = urlParts[len(urlParts)-2]
	}
	imageAndTag := urlParts[len(urlParts)-1]
	imageParts := strings.Split(imageAndTag, ":")
	image = imageParts[0]
	if len(imageParts) > 1 {
		imageTag = imageParts[len(imageParts)-1]
	}
	return image, imageTag, imageContext
}

// MajorMinorMicro ...
func MajorMinorMicro(productVersion string) (major, minor, micro string) {
	version := strings.Split(productVersion, ".")
	for len(version) < 3 {
		version = append(version, "0")
	}
	return version[0], version[1], version[2]
}
