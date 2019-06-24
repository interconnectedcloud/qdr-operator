// Package e2e defines the main entry point for qdr-operator end-to-end test suite
package e2e

import (
	"fmt"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis"
	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

// TestMain is the main entry point for the E2E test framework
func TestMain(m *testing.M) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	framework.MainEntry(m)
}

// TestQdrOperator Initializes the TestSuiteContext and
// triggers parallel execution of test Specs for the test suite.
func TestQdrOperator(t *testing.T) {

	// Initialize TestSuiteContext
	TestSuite = TestSuiteContext{
		framework.NewTestCtx(t),
		framework.Global,
	}

	// Once context is defined, set Specs to run in parallel
	testingT = t
	t.Parallel()

	// Preparing to run test suite
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	ginkgo.RunSpecsWithDefaultAndCustomReporters(
		t,
		"QDR Operator E2E Testsuite",
		[]ginkgo.Reporter{junitReporter})

}

// Initialize the framework and test context before any Spec is executed
var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {

	// TODO Use log
	interconnectList := &v1alpha1.InterconnectList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Interconnect",
			APIVersion: "interconnectedcloud.github.io/v1alpha1",
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, interconnectList)
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("failed to add custom resource scheme to framework: %v", err))
	}

	err = TestSuite.Ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: TestSuite.Ctx, Timeout: CleanupTimeout, RetryInterval: CleanupRetryInterval})
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("failed to initialize cluster resources: %v", err))
	}

	//t.Log("Initialized cluster resources")
	namespace := TestSuite.GetNamespace()

	// wait for qdr-operator to be ready
	err = e2eutil.WaitForDeployment(testingT, TestSuite.F.KubeClient, namespace, "qdr-operator", 1, RetryInterval, Timeout)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return nil
}, func(b []byte) {})

// Tear down the test and cleans up the cluster
var _ = ginkgo.SynchronizedAfterSuite(func() {
	TestSuite.Ctx.Cleanup()
}, func() {})
