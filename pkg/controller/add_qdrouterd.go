package controller

import (
	"github.com/interconnectedcloud/qdrouterd-operator/pkg/controller/qdrouterd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, qdrouterd.Add)
}
