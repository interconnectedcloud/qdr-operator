package entities

import (
	"encoding/json"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"strconv"
	"strings"
)

type ConnectionStatusType int

const (
	ConnectionStatusConnecting ConnectionStatusType = iota
	ConnectionStatusSuccess
	ConnectionStatusFailed
)

type Connector struct {
	EntityCommon
	Host                     string                         `json:"host"`
	Port                     string                         `json:"port"`
	ProtocolFamily           common.SocketAddressFamilyType `json:"protocolFamily,string"`
	Role                     common.RoleType                `json:"role,string"`
	Cost                     int                            `json:"cost"`
	SslProfile               string                         `json:"sslProfile"`
	SaslMechanisms           string                         `json:"saslMechanisms"`
	AllowRedirect            bool                           `json:"allowRedirect"`
	MaxFrameSize             int                            `json:"maxFrameSize"`
	MaxSessionFrames         int                            `json:"maxSessionFrames"`
	IdleSecondsTimeout       int                            `json:"idleSecondsTimeout"`
	StripAnnotations         StripAnnotationsType           `json:"stripAnnotations,string"`
	LinkCapacity             int                            `json:"linkCapacity"`
	VerifyHostname           bool                           `json:"verifyHostname"`
	SaslUsername             string                         `json:"saslUsername"`
	SaslPassword             string                         `json:"saslPassword"`
	MessageLoggingComponents string                         `json:"messageLoggingComponents"`
	FailoverUrls             string                         `json:"failoverUrls"`
	ConnectionStatus         ConnectionStatusType           `json:"connectionStatus,string"`
	ConnectionMsg            string                         `json:"connectionMsg"`
	PolicyVhost              string                         `json:"policyVhost"`
}

func (Connector) GetEntityId() string {
	return "connector"
}

// UnmarshalJSON returns the appropriate ConnectionStatusType for parsed string
func (a *ConnectionStatusType) UnmarshalJSON(b []byte) error {
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
	case "connecting":
		*a = ConnectionStatusConnecting
	case "success":
		*a = ConnectionStatusSuccess
	case "failed":
		*a = ConnectionStatusFailed
	}
	return nil
}

// MarshalJSON returns the string representation of ConnectionStatusType
func (a ConnectionStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	case ConnectionStatusConnecting:
		s = "CONNECTING"
	case ConnectionStatusSuccess:
		s = "SUCCESS"
	case ConnectionStatusFailed:
		s = "FAILED"
	}
	return json.Marshal(s)
}
