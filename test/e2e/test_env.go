package e2e

import (
	"fmt"
	"github.com/onsi/ginkgo"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	"testing"
	"time"
)

// TestSuiteContext Type holds the e2e Framework and TestCtx instances.
// It must be initialized before any Spec is run.
type TestSuiteContext struct {
	Ctx *framework.TestCtx
	F   *framework.Framework
}

// GetNamespace returns namespace or fail test if unable to determine it
func (ts TestSuiteContext) GetNamespace() string {
	if ts.Ctx == nil {
		ginkgo.Fail("TestSuiteContext does not contain a valid TestCtx")
	}

	namespace, err := ts.Ctx.GetNamespace()
	if err != nil {
		ginkgo.Fail(fmt.Sprintf("could not get namespace: %v", err))
	}
	return namespace
}

var (
	TestSuite TestSuiteContext
	testingT  *testing.T

	RetryInterval        = time.Second * 5
	Timeout              = time.Second * 600
	CleanupRetryInterval = time.Second * 1
	CleanupTimeout       = time.Second * 5
)
