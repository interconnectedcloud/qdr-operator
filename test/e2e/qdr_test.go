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

func TestQdr(t *testing.T) {
	qdrList := &v1alpha1.QdrList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Qdr",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, qdrList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}
	// run subtests
	t.Run("qdr-group", func(t *testing.T) {
		t.Run("Mesh", QdrCluster)
	})
}

func qdrScaleTest(t *testing.T, f *framework.Framework, ctx *framework.TestCtx) error {
	namespace, err := ctx.GetNamespace()
	if err != nil {
		return fmt.Errorf("could not get namespace: %v", err)
	}
	// create qdr customer resource
	exampleQdr := &v1alpha1.Qdr{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Qdr",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-qdr",
			Namespace: namespace,
		},
		Spec: v1alpha1.QdrSpec{
			DeploymentPlan: v1alpha1.DeploymentPlanType{
				Size:      3,
				Image:     "quay.io/interconnectedcloud/qdrouterd:1.6.0",
				Role:      "interior",
				Placement: "Any",
			},
		},
	}
	// use TestCtx's create helper to create the object and add a cleanup function for the new object
	err = f.Client.Create(goctx.TODO(), exampleQdr, &framework.CleanupOptions{TestContext: ctx, Timeout: cleanupTimeout, RetryInterval: cleanupRetryInterval})
	if err != nil {
		return err
	}
	// wait for example-qdr to reach 3 replicas
	err = e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-qdr", 3, retryInterval, timeout)
	if err != nil {
		return err
	}

	err = f.Client.Get(goctx.TODO(), types.NamespacedName{Name: "example-qdr", Namespace: namespace}, exampleQdr)
	if err != nil {
		return err
	}
	exampleQdr.Spec.DeploymentPlan.Size = 4
	err = f.Client.Update(goctx.TODO(), exampleQdr)
	if err != nil {
		return err
	}

	// wait for example-qdr to reach 4 replicas
	return e2eutil.WaitForDeployment(t, f.KubeClient, namespace, "example-qdr", 4, retryInterval, timeout)
}

func QdrCluster(t *testing.T) {
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

	if err = qdrScaleTest(t, f, ctx); err != nil {
		t.Fatal(err)
	}
}
