package integration_test

import (
	integration "github.com/cloudfoundry/bosh-bootloader/integration-test"
	"github.com/cloudfoundry/bosh-bootloader/integration-test/actors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("idempotent test", func() {
	var (
		bbl actors.BBL
	)

	BeforeEach(func() {
		var err error
		configuration, err := integration.LoadGCPConfig()
		Expect(err).NotTo(HaveOccurred())

		bbl = actors.NewBBL(configuration.StateFileDir, pathToBBL, configuration, "bbl-ci-reentrant-env")
	})

	It("is able to bbl up idempotently", func() {
		bbl.Up(actors.GCPIAAS, true)

		bbl.Up(actors.GCPIAAS, false)

		bbl.Destroy()
	})
})