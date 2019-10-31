package entities

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Address represents a Dispatch Router address
type Address struct {
	EntityCommon
	Prefix         string           `json:"prefix"`
	Pattern        string           `json:"pattern"`
	Distribution   DistributionType `json:"distribution,string"`
	Waypoint       bool             `json:"waypoint"`
	IngressPhase   int              `json:"ingressPhase"`
	EgressPhase    int              `json:"egressPhase"`
	Priority       int              `json:"priority"`
	EnableFallback bool             `json:"enableFallback"`
}

func (Address) GetEntityId() string {
	return "router.config.address"
}

type DistributionType int

const (
	DistributionMulticast DistributionType = iota
	DistributionClosest
	DistributionBalanced
	DistributionUnavailable
)

// Unmarshal JSON return the appropriate DistributionType for parsed string
func (d *DistributionType) UnmarshalJSON(b []byte) error {
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
	case "multicast":
		*d = DistributionMulticast
	case "closest":
		*d = DistributionClosest
	case "balanced":
		*d = DistributionBalanced
	case "unavailable":
		*d = DistributionUnavailable
	}
	return nil

}

// MarshalJSON returns the string represenation of DistributionType
func (d DistributionType) MarshalJSON() ([]byte, error) {
	var s string
	switch d {
	case DistributionMulticast:
		s = "multicast"
	case DistributionClosest:
		s = "closest"
	case DistributionBalanced:
		s = "balanced"
	case DistributionUnavailable:
		s = "unavailable"
	}
	return json.Marshal(s)
}
