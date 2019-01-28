package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// QdrouterdSpec defines the desired state of Qdrouterd
type QdrouterdSpec struct {
	Count                 int32        `json:"count,omitempty"`
	Image                 string       `json:"image"`
	Listeners             []Listener   `json:"listeners,:omitempty"`
	InterRouterListeners  []Listener   `json:"interRouterListeners,:omitempty"`
	SslProfiles           []SslProfile `json:"sslProfiles,omitempty"`
	Addresses             []Address    `json:"addresses,omitempty"`
	AutoLinks             []Address    `json:"autoLinks,omitempty"`
	LinkRoutes            []LinkRoute  `json:"linkRoutes,omitempty"`
	Connectors            []Connector  `json:"connectors,omitempty"`
	InterRouterConnectors []Connector  `json:"interRouterConnectors,omitempty"`
}

// QdrouterdStatus defines the observed state of Qdrouterd
type QdrouterdStatus struct {
	PodNames []string `json:"pods"`
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

type Address struct {
	Prefix       string `json:"prefix,omitempty"`
	Pattern      string `json:"pattern,omitempty"`
	Distribution string `json:"distribution,omitempty"`
	Waypoint     bool   `json:"waypoint,omitempty"`
	IngressPhase *int32 `json:"ingressPhase,omitempty"`
	EgressPhase  *int32 `json:"ingressPhase,omitempty"`
}

type Listener struct {
	Name           string `json:"name,omitempty"`
	Host           string `json:"host,omitempty"`
	Port           int32  `json:"port"`
	RouteContainer bool   `json:"role,omitempty"`
	Http           bool   `json:"http,omitempty"`
	Cost           int32  `json:"cost,omitempty"`
	SslProfile     string `json:"sslProfile,omitempty"`
}

type SslProfile struct {
	Name               string `json:"name,omitempty"`
	Credentials        string `json:"credentials,omitempty"`
	CaCert             string `json:"caCert,omitempty"`
	RequireClientCerts bool   `json:"requireClientCerts,omitempty"`
	Ciphers            string `json:"ciphers,omitempty"`
	Protocols          string `json:"protocols,omitempty"`
}

type LinkRoute struct {
	Prefix               string `json:"prefix,omitempty"`
	Pattern              string `json:"pattern,omitempty"`
	Direction            string `json:"direction,omitempty"`
	ContainerId          string `json:"containerId,omitempty"`
	Connection           string `json:"connection,omitempty"`
	AddExternalPrefix    string `json:"addExternalPrefix,omitempty"`
	RemoveExternalPrefix string `json:"removeExternalPrefix,omitempty"`
}

type Connector struct {
	Name           string `json:"name,omitempty"`
	Host           string `json:"host"`
	Port           int32  `json:"port"`
	RouteContainer bool   `json:"routeContainer,omitempty"`
	Cost           int32  `json:"cost,omitempty"`
	SslProfile     string `json:"sslProfile,omitempty"`
}

type AutoLink struct {
	Address        string `json:"address"`
	Direction      string `json:"direction"`
	ContainerId    string `json:"containerId,omitempty"`
	Connection     string `json:"connection,omitempty"`
	ExternalPrefix string `json:"externalPrefix,omitempty"`
	Phase          *int32 `json:"phase,omitempty"`
}
