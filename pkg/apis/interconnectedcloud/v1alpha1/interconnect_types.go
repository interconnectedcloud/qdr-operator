package v1alpha1

import (
	"github.com/RHsyseng/operator-utils/pkg/olm"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InterconnectSpec defines the desired state of Interconnect
type InterconnectSpec struct {
	DeploymentPlan        DeploymentPlanType `json:"deploymentPlan,omitempty"`
	Users                 string             `json:"users,omitempty"`
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
	InterconnectPhaseNone     PhaseType = ""
	InterconnectPhaseCreating           = "Creating"
	InterconnectPhaseRunning            = "Running"
	InterconnectPhaseFailed             = "Failed"
)

type ConditionType string

const (
	InterconnectConditionProvisioning ConditionType = "Provisioning"
	InterconnectConditionDeployed     ConditionType = "Deployed"
	InterconnectConditionScalingUp    ConditionType = "ScalingUp"
	InterconnectConditionScalingDown  ConditionType = "ScalingDown"
	InterconnectConditionUpgrading    ConditionType = "Upgrading"
)

type InterconnectCondition struct {
	Type           ConditionType `json:"type"`
	TransitionTime metav1.Time   `json:"transitionTime,omitempty"`
	Reason         string        `json:"reason,omitempty"`
}

// InterconnectStatus defines the observed state of Interconnect
type InterconnectStatus struct {
	Phase     PhaseType            `json:"phase,omitempty"`
	RevNumber string               `json:"revNumber,omitempty"`
	PodStatus olm.DeploymentStatus `json:"pods"`

	// Conditions keeps most recent interconnect conditions
	Conditions []InterconnectCondition `json:"conditions"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Interconnect is the Schema for the interconnects API
// +k8s:openapi-gen=true
type Interconnect struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InterconnectSpec   `json:"spec,omitempty"`
	Status InterconnectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InterconnectList contains a list of Interconnect
type InterconnectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Interconnect `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Interconnect{}, &InterconnectList{})
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
	Image        string                      `json:"image,omitempty"`
	Size         int32                       `json:"size,omitempty"`
	Role         RouterRoleType              `json:"role,omitempty"`
	Placement    PlacementType               `json:"placement,omitempty"`
	Resources    corev1.ResourceRequirements `json:"resources,omitempty"`
	Issuer       string                      `json:"issuer,omitempty"`
	LivenessPort int32                       `json:"livenessPort,omitempty"`
	ServiceType  string                      `json:"serviceType,omitempty"`
}

type Address struct {
	Prefix         string `json:"prefix,omitempty"`
	Pattern        string `json:"pattern,omitempty"`
	Distribution   string `json:"distribution,omitempty"`
	Waypoint       bool   `json:"waypoint,omitempty"`
	IngressPhase   *int32 `json:"ingressPhase,omitempty"`
	EgressPhase    *int32 `json:"egressPhase,omitempty"`
	Priority       *int32 `json:"priority,omitempty"`
	EnableFallback bool   `json:"enableFallback,omitempty"`
}

type Listener struct {
	Name             string `json:"name,omitempty"`
	Host             string `json:"host,omitempty"`
	Port             int32  `json:"port"`
	RouteContainer   bool   `json:"routeContainer,omitempty"`
	Http             bool   `json:"http,omitempty"`
	Cost             int32  `json:"cost,omitempty"`
	SslProfile       string `json:"sslProfile,omitempty"`
	SaslMechanisms   string `json:"saslMechanisms,omitempty"`
	AuthenticatePeer bool   `json:"authenticatePeer,omitempty"`
	Expose           bool   `json:"expose,omitempty"`
}

type SslProfile struct {
	Name                string `json:"name,omitempty"`
	Credentials         string `json:"credentials,omitempty"`
	CaCert              string `json:"caCert,omitempty"`
	GenerateCredentials bool   `json:"generateCredentials,omitempty"`
	GenerateCaCert      bool   `json:"generateCaCert,omitempty"`
	MutualAuth          bool   `json:"mutualAuth,omitempty"`
	Ciphers             string `json:"ciphers,omitempty"`
	Protocols           string `json:"protocols,omitempty"`
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
	VerifyHostname bool   `json:"verifyHostname,omitempty"`
	SslProfile     string `json:"sslProfile,omitempty"`
}

type AutoLink struct {
	Address        string `json:"address"`
	Direction      string `json:"direction"`
	ContainerId    string `json:"containerId,omitempty"`
	Connection     string `json:"connection,omitempty"`
	ExternalPrefix string `json:"externalPrefix,omitempty"`
	Phase          *int32 `json:"phase,omitempty"`
	Fallback       bool   `json:"fallback,omitempty"`
}
