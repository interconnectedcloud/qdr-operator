package entities

import (
	"encoding/json"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"strconv"
	"strings"
)

type AutoLink struct {
	EntityCommon
	Address         string                 `json:"address"`
	Direction       common.DirectionType   `json:"direction,string"`
	Phase           int                    `json:"phase"`
	ContainerId     string                 `json:"containerId"`
	Connection      string                 `json:"connection"`
	ExternalAddress string                 `json:"externalAddress"`
	Fallback        bool                   `json:"fallback"`
	LinkRef         string                 `json:"linkRef"`
	OperStatus      AutoLinkOperStatusType `json:"operStatus,string"`
	LastError       string                 `json:"lastError"`
}

func (AutoLink) GetEntityId() string {
	return "router.config.autoLink"
}

type AutoLinkOperStatusType int

const (
	AutoLinkOperStatusInactive AutoLinkOperStatusType = iota
	AutoLinkOperStatusAttaching
	AutoLinkOperStatusFailed
	AutoLinkOperStatusActive
	AutoLinkOperStatusQuiescing
	AutoLinkOperStatusIdle
)

// UnmarshalJSON returns the appropriate AutoLinkOperStatusType for parsed string
func (o *AutoLinkOperStatusType) UnmarshalJSON(b []byte) error {
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
	case "inactive":
		*o = AutoLinkOperStatusInactive
	case "attaching":
		*o = AutoLinkOperStatusAttaching
	case "failed":
		*o = AutoLinkOperStatusFailed
	case "active":
		*o = AutoLinkOperStatusActive
	case "quiescing":
		*o = AutoLinkOperStatusQuiescing
	case "idle":
		*o = AutoLinkOperStatusIdle
	}
	return nil
}

// MarshalJSON returns the string representation of AutoLinkOperStatusType
func (o AutoLinkOperStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch o {
	case AutoLinkOperStatusInactive:
		s = "inactive"
	case AutoLinkOperStatusAttaching:
		s = "attaching"
	case AutoLinkOperStatusFailed:
		s = "failed"
	case AutoLinkOperStatusActive:
		s = "active"
	case AutoLinkOperStatusQuiescing:
		s = "quiescing"
	case AutoLinkOperStatusIdle:
		s = "idle"
	}
	return json.Marshal(s)
}
