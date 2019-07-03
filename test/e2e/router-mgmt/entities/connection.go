package entities

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Connection represents a Dispatch Router connection
type Connection struct {
	EntityCommon
	Active          bool                   `json:"active"`
	AdminStatus     AdminStatusType        `json:"adminStatus,string"`
	OperStatus      OperStatusType         `json:"operStatus,string"`
	Container       string                 `json:"container"`
	Opened          bool                   `json:"opened"`
	Host            string                 `json:"host"`
	Direction       DirectionType          `json:"dir,string"`
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

type OperStatusType int

const (
	OperStatusUp OperStatusType = iota
	OperStatusClosing
)

// UnmarshalJSON returns the appropriate OperStatusType for parsed string
func (o *OperStatusType) UnmarshalJSON(b []byte) error {
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
		*o = OperStatusUp
	case "closing":
		*o = OperStatusClosing
	}
	return nil
}

// MarshalJSON returns the string representation of OperStatusType
func (o OperStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch o {
	case OperStatusUp:
		s = "up"
	case OperStatusClosing:
		s = "closing"
	}
	return json.Marshal(s)
}

type DirectionType int

const (
	DirectionTypeIn DirectionType = iota
	DirectionTypeOut
)

// UnmarshalJSON returns the appropriate DirectionType for parsed string
func (d *DirectionType) UnmarshalJSON(b []byte) error {
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
		*d = DirectionTypeIn
	case "out":
		*d = DirectionTypeOut
	}
	return nil
}

// MarshalJSON returns the string representation of DirectionType
func (d DirectionType) MarshalJSON() ([]byte, error) {
	var s string
	switch d {
	case DirectionTypeIn:
		s = "in"
	case DirectionTypeOut:
		s = "out"
	}
	return json.Marshal(s)
}
