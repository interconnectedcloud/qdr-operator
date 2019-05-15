package e2e

import (
	goctx "context"
	"fmt"
	"testing"
	"time"

	apis "github.com/interconnectedcloud/qdr-operator/pkg/apis"
	v1alpha1 "github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"

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

func TestInterconnect(t *testing.T) {
	interconnectList := &v1alpha1.InterconnectList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Interconnect",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, interconnectList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("interconnect-group", func(t *testing.T) {
		t.Run("Mesh", InterconnectCluster)
	})
}

func interconnectScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create interconnect customer resource
	exampleInterconnect := &v1alpha1.Interconnect{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Interconnect",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-interconnect",
			Namespace: namespace,
		},
		Spec: v1alpha1.InterconnectSpec{
			DeploymentPlan: v1alpha1.DeploymentPlanType{
				Size:      3,
				Image:     "quay.io/interconnectedcloud/qdrouterd:1.6.0",
				Role:      "interior",
				Placement: "Any",
			},
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleInterconnect, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-interconnect to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-interconnect", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-interconnect", Namespace: namespace}, exampleInterconnect)
	if err != nil {
		return err
	}
	exampleInterconnect.Spec.DeploymentPlan.Size = 4
	err = f.Client.Update(goctx.TODO(), exampleInterconnect)
	if err != nil {
		return err
	}

	// wait for example-interconnect to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-interconnect", 4, retryInterval, timeout)
}

func InterconnectCluster(t *testing.T) {
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
	// wait for qdr-operator to be ready
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "qdr-operator", 1, retryInterval, timeout)
	if err != nil {
		t.Fatal(err)
	}

	if err = interconnectScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
