package client

import (
	"context"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e"
	"github.com/onsi/ginkgo"
	"github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

// DeployAndValidate deploys the given cr and validate if the client
// is done and the number of replicas match what is defined in the CR.
func DeployAndValidate(cr *v1alpha1.Interconnect, t *testing.T) {

	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err := e2e.TestSuite.F.Client.Create(context.TODO(), cr, &test.CleanupOptions{
		TestContext:   e2e.TestSuite.Ctx,
		Timeout:       e2e.CleanupTimeout,
		RetryInterval: e2e.CleanupRetryInterval,
	})
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	// wait for example-interconnect to reach number of replicas defined in DeploymentPlan
	err = e2eutil.WaitForDeployment(t, e2e.TestSuite.F.KubeClient, cr.ObjectMeta.Namespace,
		cr.ObjectMeta.Name, int(cr.Spec.DeploymentPlan.Size), e2e.RetryInterval, e2e.Timeout)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	// validate the client is present in the namespace
	err = e2e.TestSuite.F.Client.Get(context.TODO(),
		types.NamespacedName{Name: cr.ObjectMeta.Name, Namespace: cr.ObjectMeta.Namespace}, cr)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

}

// UpdateDeploymentPlanSize will update a previously deployed CR
// and then validate if deployment succeeded and number of replicas
// match what is defined in the CR.
func UpdateDeploymentPlanSize(cr *v1alpha1.Interconnect, t *testing.T) {
	err := e2e.TestSuite.F.Client.Update(context.TODO(), cr)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	// wait for example-interconnect to reach # of replicas given
	err = e2eutil.WaitForDeployment(t, e2e.TestSuite.F.KubeClient, cr.ObjectMeta.Namespace,
		cr.ObjectMeta.Name, int(cr.Spec.DeploymentPlan.Size), e2e.RetryInterval, e2e.Timeout)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	// validate the client is present in the namespace
	err = e2e.TestSuite.F.Client.Get(context.TODO(),
		types.NamespacedName{Name: cr.ObjectMeta.Name, Namespace: cr.ObjectMeta.Namespace}, cr)
	if err != nil {
		ginkgo.Fail(err.Error())
	}
}
