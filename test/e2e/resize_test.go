package e2e

import (
	"context"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework/qdrmanagement"

	"github.com/interconnectedcloud/qdr-operator/pkg/apis/interconnectedcloud/v1alpha1"
	"github.com/interconnectedcloud/qdr-operator/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("[resize_test] Interconnect resize deployment tests", func() {
	f := framework.NewFramework("basic-interior", nil)

	It("Should be able to resize an interior deployment from 3 to 5", func() {
		testInteriorResize(f, 3, 5)
	})

	It("Should be able to resize an interior deployment from 5 to 3", func() {
		testInteriorResize(f, 5, 3)
	})
})

func testInteriorResize(f *framework.Framework, initialSize int, finalSize int) {
	var (
		name = "interior-interconnect"
	)

	By("Creating an interior interconnect with initial size")
	ei, err := f.CreateInterconnect(f.Namespace, int32(initialSize), func(ei *v1alpha1.Interconnect) {
		ei.Name = name
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

	By("Waiting until full interconnect with initial size and version")
	ctx1, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = f.WaitUntilFullInterconnectWithSize(ctx1, ei, initialSize)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect initial qdr entities")
	ctx2, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx2, f, ei)
	Expect(err).NotTo(HaveOccurred())

	By("Retrieving the Interconnect resource in the namespace")
	ei, err = f.GetInterconnect(name)
	Expect(err).NotTo(HaveOccurred())

	By("Scaling the interior interconnect to final size")
	ei.Spec.DeploymentPlan.Size = int32(finalSize)
	_, err = f.UpdateInterconnect(ei)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect with final size and version")
	ctx3, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = f.WaitUntilFullInterconnectWithSize(ctx3, ei, finalSize)
	Expect(err).NotTo(HaveOccurred())

	By("Waiting until full interconnect with final qdr entities")
	ctx4, fn := context.WithTimeout(context.Background(), framework.Timeout)
	defer fn()
	err = qdrmanagement.WaitUntilFullInterconnectWithQdrEntities(ctx4, f, ei)
	Expect(err).NotTo(HaveOccurred())

}
