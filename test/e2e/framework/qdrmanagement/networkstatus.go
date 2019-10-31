package qdrmanagement

import (
	"context"

	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement/entities"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// WaitForQdrNodesInPod attempts to retrieve the list of Node Entities
// present on the given pod till the expected amount of nodes are present
// or an error or timeout occurs.
func WaitForQdrNodesInPod(f *framework.Framework, pod v1.Pod, expected int, retryInterval, timeout time.Duration) error {
	// Retry logic to retrieve nodes
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		nodes, err := QdmanageQuery(f, pod.Name, entities.Node{}, nil)
		if err != nil {
			return false, err
		}
		if len(nodes) != expected {
			return false, nil
		}
		return true, nil
	})
	return err
}

// InterRouterConnectionsForPod will get all opened inter-router connections
func InterRouterConnectionsForPod(f *framework.Framework, pod v1.Pod) ([]entities.Connection, error) {
	conns, err := QdmanageQuery(f, pod.Name, entities.Connection{}, func(entity entities.Entity) bool {
		conn := entity.(entities.Connection)
		if conn.Role == "inter-router" && conn.Opened {
			return true
		}
		return false
	})

	// Preparing result
	connections := make([]entities.Connection, 0)
	for _, c := range conns {
		connections = append(connections, c.(entities.Connection))
	}
	return connections, err
}

func InterconnectHasExpectedNodes(f *framework.Framework, interconnect *v1alpha1.Interconnect) (bool, error) {
	pods, err := f.PodsForInterconnect(interconnect)
	if err != nil {
		return false, err
	}

	for _, pod := range pods {
		nodes, err := QdmanageQuery(f, pod.Name, entities.Node{}, nil)
		if err != nil {
			return false, err
		}
		if int32(len(nodes)) != interconnect.Spec.DeploymentPlan.Size {
			return false, nil
		}
	}
	return true, nil
}

func InterconnectHasExpectedInterRouterConnections(f *framework.Framework, interconnect *v1alpha1.Interconnect) (bool, error) {
	if interconnect.Spec.DeploymentPlan.Role != v1alpha1.RouterRoleInterior {
		// edge role, nothing to see here
		return true, nil
	}

	pods, err := f.PodsForInterconnect(interconnect)
	if err != nil {
		return false, err
	}

	for _, pod := range pods {
		nodes, err := QdmanageQuery(f, pod.Name, entities.Node{}, nil)
		if err != nil {
			return false, err
		}
		if int32(len(nodes)) != interconnect.Spec.DeploymentPlan.Size-1 {
			return false, nil
		}
	}
	return true, nil
}

// Wait until all the pods belonging to the Interconnect deployment report
// expected node counts, irc's, etc.
func WaitUntilFullInterconnectWithQdrEntities(ctx context.Context, f *framework.Framework, interconnect *v1alpha1.Interconnect) error {

	return framework.RetryWithContext(ctx, framework.RetryInterval, func() (bool, error) {
		// Check that all the qdr pods have the expected node cound
		n, err := InterconnectHasExpectedNodes(f, interconnect)
		if err != nil {
			return false, nil
		}
		if !n {
			return false, nil
		}

		// Check that all the qdr pods have the expected inter router connections
		//i, err := InterconnectHasExpectedInterRouterConnections(f, interconnect)
		//if err != nil {
		//    return false, nil
		//}
		//if !i {
		//    return false, nil
		//}
		return true, nil
	})
}
