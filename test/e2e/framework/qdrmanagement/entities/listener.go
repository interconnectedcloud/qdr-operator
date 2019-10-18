package entities

import (
	"encoding/json"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"strconv"
	"strings"
)

type Listener struct {
	EntityCommon
	Host                           string                         `json:"host"`
	Port                           string                         `json:"port"`
	SocketAddressFamily            common.SocketAddressFamilyType `json:"socketAddressFamily,string"`
	Role                           common.RoleType                `json:"role,string"`
	Cost                           int                            `json:"cost"`
	SslProfile                     string                         `json:"sslProfile"`
	SaslMechanisms                 string                         `json:"saslMechanisms"`
	AuthenticatePeer               bool                           `json:"authenticatePeer"`
	SaslPlugin                     string                         `json:"saslPlugin"`
	RequireEncryption              bool                           `json:"requireEncryption"`
	RequireSsl                     bool                           `json:"requireSsl"`
	TrustedCertsFile               string                         `json:"trustedCertsFile"`
	MaxFrameSize                   int                            `json:"maxFrameSize"`
	MaxSessions                    int                            `json:"maxSessions"`
	MaxSessionFrames               int                            `json:"maxSessionFrames"`
	IdleTimeoutSeconds             int                            `json:"idleTimeoutSeconds"`
	InitialHandshakeTimeoutSeconds int                            `json:"initialHandshakeTimeoutSeconds"`
	StripAnnotations               StripAnnotationsType           `json:"stripAnnotations,string"`
	LinkCapacity                   int                            `json:"linkCapacity"`
	MultiTenant                    bool                           `json:"multiTenant"`
	FailoverUrls                   string                         `json:"failoverUrls"`
	Healthz                        bool                           `json:"healthz"`
	Metrics                        bool                           `json:"metrics"`
	Websockets                     bool                           `json:"websockets"`
	Http                           bool                           `json:"http"`
	HttpRootDir                    string                         `json:"httpRootDir"`
	MessageLoggingComponents       string                         `json:"messageLoggingComponents"`
	PolicyVhost                    string                         `json:"policyVhost"`
}

func (Listener) GetEntityId() string {
	return "listener"
}

type StripAnnotationsType int

const (
	StripAnnotationsTypeIn StripAnnotationsType = iota
	StripAnnotationsTypeOut
	StripAnnotationsTypeBoth
	StripAnnotationsTypeNo
)

// UnmarshalJSON returns the appropriate StripAnnotationsType for parsed string
func (d *StripAnnotationsType) UnmarshalJSON(b []byte) error {
	var s string

	if len(b) == 0 {
		return nil
	}
	if b[0] != '"' {
		b = []byte(strconv.Quote(string(b)))
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "in":
		*d = StripAnnotationsTypeIn
	case "out":
		*d = StripAnnotationsTypeOut
	case "no":
		*d = StripAnnotationsTypeNo
	default:
		*d = StripAnnotationsTypeBoth
	}
	return nil
}

// MarshalJSON returns the string representation of StripAnnotationsType
func (d StripAnnotationsType) MarshalJSON() ([]byte, error) {
	var s string
	switch d {
	case StripAnnotationsTypeIn:
		s = "in"
	case StripAnnotationsTypeOut:
		s = "out"
	case StripAnnotationsTypeNo:
		s = "no"
	default:
		s = "both"
	}
	return json.Marshal(s)
}
