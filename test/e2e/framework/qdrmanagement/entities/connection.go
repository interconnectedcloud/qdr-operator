package entities

import (
	"encoding/json"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"strconv"
	"strings"
)

// Connection represents a Dispatch Router connection
type Connection struct {
	EntityCommon
	Active          bool                   `json:"active"`
	AdminStatus     AdminStatusType        `json:"adminStatus,string"`
	OperStatus      ConnOperStatusType     `json:"operStatus,string"`
	Container       string                 `json:"container"`
	Opened          bool                   `json:"opened"`
	Host            string                 `json:"host"`
	Direction       common.DirectionType   `json:"dir,string"`
	Role            string                 `json:"role"`
	IsAuthenticated bool                   `json:"isAuthenticated"`
	IsEncrypted     bool                   `json:"isEncrypted"`
	Sasl            string                 `json:"sasl"`
	User            string                 `json:"user"`
	Ssl             bool                   `json:"ssl"`
	SslProto        string                 `json:"sslProto"`
	SslCipher       string                 `json:"sslCipher"`
	SslSsf          int                    `json:"sslSsf"`
	Tenant          string                 `json:"tenant"`
	Properties      map[string]interface{} `json:"properties"`
}

// Implementation of Entity interface
func (Connection) GetEntityId() string {
	return "connection"
}

type AdminStatusType int

const (
	AdminStatusEnabled AdminStatusType = iota
	AdminStatusDeleted
)

// UnmarshalJSON returns the appropriate AdminStatusType for parsed string
func (a *AdminStatusType) UnmarshalJSON(b []byte) error {
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
	case "enabled":
		*a = AdminStatusEnabled
	case "deleted":
		*a = AdminStatusDeleted
	}
	return nil
}

// MarshalJSON returns the string representation of AdminStatusType
func (a AdminStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	case AdminStatusEnabled:
		s = "enabled"
	case AdminStatusDeleted:
		s = "deleted"
	}
	return json.Marshal(s)
}

type ConnOperStatusType int

const (
	ConnOperStatusUp ConnOperStatusType = iota
	ConnOperStatusClosing
)

// UnmarshalJSON returns the appropriate ConnOperStatusType for parsed string
func (o *ConnOperStatusType) UnmarshalJSON(b []byte) error {
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
	case "up":
		*o = ConnOperStatusUp
	case "closing":
		*o = ConnOperStatusClosing
	}
	return nil
}

// MarshalJSON returns the string representation of ConnOperStatusType
func (o ConnOperStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch o {
	case ConnOperStatusUp:
		s = "up"
	case ConnOperStatusClosing:
		s = "closing"
	}
	return json.Marshal(s)
}
