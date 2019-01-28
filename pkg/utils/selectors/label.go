package selectors

import (
	"k8s.io/apimachinery/pkg/labels"
)

const (
	LabelAppKey = "application"

	LabelResourceKey = "qdrouterd_cr"
)

// Set labels in a map
func LabelsForQdrouterd(name string) map[string]string {
	return map[string]string{
		LabelAppKey:      name,
		LabelResourceKey: name,
	}
}

// return a selector that matches resources for a qdrouterd resource
func ResourcesByQdrouterdName(name string) labels.Selector {
	set := map[string]string{
		LabelAppKey:      name,
		LabelResourceKey: name,
	}
	return labels.SelectorFromSet(set)
}
