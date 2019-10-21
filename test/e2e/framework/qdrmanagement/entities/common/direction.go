package common

import (
	"encoding/json"
	"strconv"
	"strings"
)

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
