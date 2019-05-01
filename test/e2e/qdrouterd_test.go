package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	apis "github.com/interconnectedcloud/qdrouterd-operator/pkg/apis"
	v1alpha1 "github.com/interconnectedcloud/qdrouterd-operator/pkg/apis/interconnectedcloud/v1alpha1"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	retryInterval        = time.Second * 5
	timeout              = time.Second * 600
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestQdrouterd(t *testing.T) {
	qdrouterdList := &v1alpha1.QdrouterdList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Qdrouterd",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, qdrouterdList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("qdrouterd-group", func(t *testing.T) {
		t.Run("Mesh", QdrouterdCluster)
	})
}

func qdrouterdScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create qdrouterd customer resource
	exampleQdrouterd := &v1alpha1.Qdrouterd{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Qdrouterd",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-qdrouterd",
			Namespace: namespace,
		},
		Spec: v1alpha1.QdrouterdSpec{
			DeploymentPlan: v1alpha1.DeploymentPlanType{
				Size:      3,
				Image:     "quay.io/ajssmith/qpid-dispatch-router:1.6.0",
				Role:      "interior",
				Placement: "Any",
			},
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleQdrouterd, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-qdrouterd to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-qdrouterd", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-qdrouterd", Namespace: namespace}, exampleQdrouterd)
	if err != nil {
		return err
	}
	exampleQdrouterd.Spec.DeploymentPlan.Size = 4
	err = f.Client.Update(goctx.TODO(), exampleQdrouterd)
	if err != nil {
		return err
	}

	// wait for example-qdrouterd to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-qdrouterd", 4, retryInterval, timeout)
}

func QdrouterdCluster(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	err := ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}
	t.Log("Initialized cluster resources")
	namespace, err := ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework variables
	f := framework.Global
	// wait for qdrouterd-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "qdrouterd-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = qdrouterdScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
