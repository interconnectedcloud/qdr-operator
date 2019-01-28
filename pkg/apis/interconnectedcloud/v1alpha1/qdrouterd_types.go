package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// QdrouterdSpec defines the desired state of Qdrouterd
type QdrouterdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// QdrouterdStatus defines the observed state of Qdrouterd
type QdrouterdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Qdrouterd is the Schema for the qdrouterds API
// +k8s:openapi-gen=true
type Qdrouterd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrouterdSpec   `json:"spec,omitempty"`
	Status QdrouterdStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QdrouterdList contains a list of Qdrouterd
type QdrouterdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Qdrouterd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Qdrouterd{}, &QdrouterdList{})
}
