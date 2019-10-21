package common

import (
	"encoding/json"
	"strconv"
	"strings"
)

// SocketAddressFamilyType
type SocketAddressFamilyType int

const (
	SocketAddressFamilyIPv4 SocketAddressFamilyType = iota
	SocketAddressFamilyIPv6
)

// UnmarshalJSON returns the appropriate SocketAddressFamilyType for parsed string
func (a *SocketAddressFamilyType) UnmarshalJSON(b []byte) error {
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
	case "IPv6":
		*a = SocketAddressFamilyIPv6
	default:
		*a = SocketAddressFamilyIPv4
	}
	return nil
}

// MarshalJSON returns the string representation of SocketAddressFamilyType
func (a SocketAddressFamilyType) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	case SocketAddressFamilyIPv6:
		s = "IPv6"
	default:
		s = "IPv4"
	}
	return json.Marshal(s)
}
