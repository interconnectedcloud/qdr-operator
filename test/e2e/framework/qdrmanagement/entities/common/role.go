package common

import (
	"encoding/json"
	"strconv"
	"strings"
)

// RoleType
type RoleType int

const (
	RoleNormal RoleType = iota
	RoleInterRouter
	RoleRouteContainer
	RoleEdge
)

// UnmarshalJSON returns the appropriate RoleType for parsed string
func (a *RoleType) UnmarshalJSON(b []byte) error {
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
	case "inter-router":
		*a = RoleInterRouter
	case "route-container":
		*a = RoleRouteContainer
	case "edge":
		*a = RoleEdge
	default:
		*a = RoleNormal
	}
	return nil
}

// MarshalJSON returns the string representation of RoleType
func (a RoleType) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	case RoleInterRouter:
		s = "inter-router"
	case RoleRouteContainer:
		s = "route-container"
	case RoleEdge:
		s = "edge"
	default:
		s = "normal"
	}
	return json.Marshal(s)
}
