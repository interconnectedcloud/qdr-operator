package entities

import (
	"encoding/json"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities/common"
	"strconv"
	"strings"
)

type LinkRoute struct {
	EntityCommon
	Prefix            string                  `json:"prefix"`
	Pattern           string                  `json:"pattern"`
	AddExternalPrefix string                  `json:"addExternalPrefix"`
	DelExternalPrefix string                  `json:"delExternalPrefix"`
	ContainerId       string                  `json:"containerId"`
	Connection        string                  `json:"connection"`
	Distribution      LinkDistributionType    `json:"distribution,string"`
	Direction         common.DirectionType    `json:"direction,string"`
	OperStatus        LinkRouteOperStatusType `json:"operStatus,string"`
}

func (LinkRoute) GetEntityId() string {
	return "router.config.linkRoute"
}

type LinkDistributionType int

const (
	LinkDistributionBalanced LinkDistributionType = iota
)

// Unmarshal JSON return the appropriate LinkDistributionType for parsed string
func (d *LinkDistributionType) UnmarshalJSON(b []byte) error {
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
	case "linkBalanced":
		*d = LinkDistributionBalanced
	}
	return nil

}

// MarshalJSON returns the string represenation of LinkDistributionType
func (d LinkDistributionType) MarshalJSON() ([]byte, error) {
	var s string
	switch d {
	case LinkDistributionBalanced:
		s = "linkBalanced"
	}
	return json.Marshal(s)
}

type LinkRouteOperStatusType int

const (
	LinkRouteOperStatusInactive LinkRouteOperStatusType = iota
	LinkRouteOperStatusActive
)

// Unmarshal JSON return the appropriate LinkRouteOperStatusType for parsed string
func (d *LinkRouteOperStatusType) UnmarshalJSON(b []byte) error {
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
		*d = LinkRouteOperStatusInactive
	case "active":
		*d = LinkRouteOperStatusActive
	}
	return nil

}

// MarshalJSON returns the string represenation of LinkRouteOperStatusType
func (d LinkRouteOperStatusType) MarshalJSON() ([]byte, error) {
	var s string
	switch d {
	case LinkRouteOperStatusInactive:
		s = "inactive"
	case LinkRouteOperStatusActive:
		s = "active"
	}
	return json.Marshal(s)
}
