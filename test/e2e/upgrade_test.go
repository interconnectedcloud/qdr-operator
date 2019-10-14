package e2e

import (
	"context"

	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[upgrade_test] Interconnect upgrade deployment tests", func() {
	f := framework.NewFramework("basic-interior", nil)

	It("Should be able to upgrade the qdrouterd image for an interior deployment", func() {
		testInteriorImageUpgrade(f)
	})

})

func testInteriorImageUpgrade(f *framework.Framework) {
	var (
		name             = "interior-interconnect"
		image            = "quay.io/interconnectedcloud/qdrouterd"
		initialVersion   = "1.8.0"
		finalVersion     = "1.9.0"
		size             = 3
	)

	By("Creating an interior interconnect with size 3")
	ei, err := f.CreateInterconnect(f.Namespace, int32(size), func(ei *v1alpha1.Interconnect) {
		ei.Name = name
		ei.Spec.DeploymentPlan.Image = image + ":" + initialVersion
	})
	Expect(err).NotTo(HaveOccurred())

	// Make sure we cleanup the Interconnect resource after we're done testing.
	defer func() {
		err = f.DeleteInterconnect(ei)
		Expect(err).NotTo(HaveOccurred())
	}()

	By("Creating an Interconnect resource in the namespace")
	ei, err = f.GetInterconnect(name)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect with size and initial version")
	ctx1, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = f.WaitUntilFullInterconnectWithVersion(ctx1, ei, size, initialVersion)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect initial qdr entities")
	ctx2, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx2, f, ei)
	Expect(err).NotTo(HaveOccurred())

	By("Retrieving the Interconnect resource in the namespace")
	ei, err = f.GetInterconnect(name)
	Expect(err).NotTo(HaveOccurred())

	By("Upgrading the qdrouterd image version")
	ei.Spec.DeploymentPlan.Image = image + ":" + finalVersion
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect with size and final version")
	ctx3, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = f.WaitUntilFullInterconnectWithVersion(ctx3, ei, size, finalVersion)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect with final qdr entities")
	ctx4, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx4, f, ei)
	Expect(err).NotTo(HaveOccurred())

}
