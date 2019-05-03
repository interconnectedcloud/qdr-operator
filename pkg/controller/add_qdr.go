package controller

import (
	"github.com/interconnectedcloud/qdr-operator/pkg/controller/qdr"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, qdr.Add)
}
