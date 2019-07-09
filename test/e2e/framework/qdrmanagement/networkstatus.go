package qdrmanagement

import (
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// WaitForQdrNodesInPod attempts to retrieve the list of Node Entities
// present on the given pod till the expected amount of nodes are present
// or an error or timeout occurs.
func WaitForQdrNodesInPod(f *framework.Framework, pod v1.Pod, expected int) error {
	var nodes []entities.Node
	// Retry logic to retrieve nodes
	err := wait.Poll(5*time.Second, 20*time.Second, func() (done bool, err error) {
		if nodes, err = QdmanageQueryNodes(f, pod.Name); err != nil {
			return false, err
		}
		if len(nodes) != expected {
			return false, nil
		}
		return true, nil
	})
	return err
}

// ListInterRouterConnectionsForPod will get all opened inter-router connections
func ListInterRouterConnectionsForPod(f *framework.Framework, pod v1.Pod) ([]entities.Connection, error) {
	conns, err := QdmanageQueryConnectionsFilter(f, pod.Name, func(entity interface{}) bool {
		conn := entity.(entities.Connection)
		if conn.Role == "inter-router" && conn.Opened {
			return true
		}
		return false
	})
	return conns, err
}
