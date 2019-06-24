package configmaps

import (
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/configs"
	"github.com/interconnectedcloud/qdr-operator/pkg/utils/selectors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Create NewConfigMapForCR method to create configmap
func NewConfigMapForCR(m *v1alpha1.Interconnect) *corev1.ConfigMap {
	config := configs.ConfigForInterconnect(m)
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Data: map[string]string{
			"qdrouterd.conf.template": config,
		},
	}

	return configMap
}

func NewConfigMapForSaslConfig(m *v1alpha1.Interconnect) *corev1.ConfigMap {
	labels := selectors.LabelsForInterconnect(m.Name)
	config := configs.ConfigForSasl(m)
	configMap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    labels,
			Name:      m.Name + "-sasl-config",
			Namespace: m.Namespace,
		},
		Data: map[string]string{
			"qdrouterd.conf": config,
		},
	}

	return configMap
}
