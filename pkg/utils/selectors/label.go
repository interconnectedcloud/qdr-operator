package selectors

import (
	"k8s.io/apimachinery/pkg/labels"
)

const (
	LabelAppKey = "application"

	LabelResourceKey = "qdr_cr"
)

// Set labels in a map
func LabelsForQdr(name string) map[string]string {
	return map[string]string{
		LabelAppKey:      name,
		LabelResourceKey: name,
	}
}

// return a selector that matches resources for a qdr resource
func ResourcesByQdrName(name string) labels.Selector {
	set := map[string]string{
		LabelAppKey:      name,
		LabelResourceKey: name,
	}
	return labels.SelectorFromSet(set)
}
