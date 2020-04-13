package constants

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RedHatImageRegistry = "registry.redhat.io"
	OPeratorImagTag     = "1.2"

	InterconnectVar         = "QDROUTERD_IMAGE"
	Interconnect17Image     = "amq-interconnect"
	Interconnect17ImageTag  = "1.7"
	Interconnect17ImageURL  = RedHatImageRegistry + "/amq7/" + Interconnect17Image + ":" + Interconnect17ImageTag
	Interconnect17Component = "amq-interconnect-openshift-container"
)

type ImageEnv struct {
	Var       string
	Component string
	Registry  string
}
type ImageRef struct {
	metav1.TypeMeta `json:",inline"`
	Spec            ImageRefSpec `json:"spec"`
}
type ImageRefSpec struct {
	Tags []ImageRefTag `json:"tags"`
}
type ImageRefTag struct {
	Name string                  `json:"name"`
	From *corev1.ObjectReference `json:"from"`
}
