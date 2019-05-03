package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// QdrSpec defines the desired state of Qdr
type QdrSpec struct {
	DeploymentPlan        DeploymentPlanType `json:"deploymentPlan,omitempty"`
	Listeners             []Listener         `json:"listeners,omitempty"`
	InterRouterListeners  []Listener         `json:"interRouterListeners,omitempty"`
	EdgeListeners         []Listener         `json:"edgeListeners,omitempty"`
	SslProfiles           []SslProfile       `json:"sslProfiles,omitempty"`
	Addresses             []Address          `json:"addresses,omitempty"`
	AutoLinks             []AutoLink         `json:"autoLinks,omitempty"`
	LinkRoutes            []LinkRoute        `json:"linkRoutes,omitempty"`
	Connectors            []Connector        `json:"connectors,omitempty"`
	InterRouterConnectors []Connector        `json:"interRouterConnectors,omitempty"`
	EdgeConnectors        []Connector        `json:"edgeConnectors,omitempty"`
}

type PhaseType string

const (
	QdrPhaseNone     PhaseType = ""
	QdrPhaseCreating           = "Creating"
	QdrPhaseRunning            = "Running"
	QdrPhaseFailed             = "Failed"
)

type ConditionType string

const (
	QdrConditionProvisioning ConditionType = "Provisioning"
	QdrConditionDeployed     ConditionType = "Deployed"
	QdrConditionScalingUp    ConditionType = "ScalingUp"
	QdrConditionScalingDown  ConditionType = "ScalingDown"
	QdrConditionUpgrading    ConditionType = "Upgrading"
)

type QdrCondition struct {
	Type           ConditionType `json:"type"`
	TransitionTime metav1.Time   `json:"transitionTime,omitempty"`
	Reason         string        `json:"reason,omitempty"`
}

// QdrStatus defines the observed state of Qdr
type QdrStatus struct {
	Phase     PhaseType `json:"phase,omitempty"`
	RevNumber string    `json:"revNumber,omitempty"`
	PodNames  []string  `json:"pods"`

	// Conditions keeps most recent qdr conditions
	Conditions []QdrCondition `json:"conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Qdr is the Schema for the qdrs API
// +k8s:openapi-gen=true
type Qdr struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   QdrSpec   `json:"spec,omitempty"`
	Status QdrStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QdrList contains a list of Qdr
type QdrList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Qdr `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Qdr{}, &QdrList{})
}

type RouterRoleType string

const (
	RouterRoleInterior RouterRoleType = "interior"
	RouterRoleEdge                    = "edge"
)

type PlacementType string

const (
	PlacementAny          PlacementType = "Any"
	PlacementEvery                      = "Every"
	PlacementAntiAffinity               = "AntiAffinity"
	PlacementNode                       = "Node"
)

type DeploymentPlanType struct {
	Image     string                      `json:"image,omitempty"`
	Size      int32                       `json:"size,omitempty"`
	Role      RouterRoleType              `json:"role,omitempty"`
	Placement PlacementType               `json:"placement,omitempty"`
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	Issuer    string                      `json:"issuer,omitempty"`
}

type Address struct {
	Prefix       string `json:"prefix,omitempty"`
	Pattern      string `json:"pattern,omitempty"`
	Distribution string `json:"distribution,omitempty"`
	Waypoint     bool   `json:"waypoint,omitempty"`
	IngressPhase *int32 `json:"ingressPhase,omitempty"`
	EgressPhase  *int32 `json:"egressPhase,omitempty"`
	Priority     *int32 `json:"priority,omitempty"`
}

type Listener struct {
	Name           string `json:"name,omitempty"`
	Host           string `json:"host,omitempty"`
	Port           int32  `json:"port"`
	RouteContainer bool   `json:"routeContainer,omitempty"`
	Http           bool   `json:"http,omitempty"`
	Cost           int32  `json:"cost,omitempty"`
	SslProfile     string `json:"sslProfile,omitempty"`
	Expose         bool   `json:"expose,omitempty"`
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
