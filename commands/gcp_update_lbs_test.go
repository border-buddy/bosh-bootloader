package commands_test

import (
	"github.com/cloudfoundry/bosh-bootloader/commands"
	"github.com/cloudfoundry/bosh-bootloader/fakes"
	"github.com/cloudfoundry/bosh-bootloader/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GCP Update LBs", func() {
	var (
		command      commands.GCPUpdateLBs
		gcpCreateLBs *fakes.GCPCreateLBs
		state        storage.State
	)

	BeforeEach(func() {
		gcpCreateLBs = &fakes.GCPCreateLBs{}

		command = commands.NewGCPUpdateLBs(gcpCreateLBs)

		state = storage.State{
			IAAS: "gcp",
			LB: storage.LB{
				Type:   "cf",
				Cert:   "some-cert",
				Key:    "some-key",
				Domain: "some-domain",
			},
		}
	})

	Describe("Execute", func() {
		It("calls out to GCP Create LBs", func() {
			config := commands.GCPCreateLBsConfig{
				CertPath: "some-cert-path",
				KeyPath:  "some-key-path",
				LBType:   "cf",
				Domain:   "some-domain",
			}
			err := command.Execute(config, state)

			Expect(err).NotTo(HaveOccurred())
			Expect(gcpCreateLBs.ExecuteCall.CallCount).To(Equal(1))
			Expect(gcpCreateLBs.ExecuteCall.Receives.Config).To(Equal(commands.GCPCreateLBsConfig{
				CertPath: "some-cert-path",
				KeyPath:  "some-key-path",
				LBType:   "cf",
				Domain:   "some-domain",
			}))
			Expect(gcpCreateLBs.ExecuteCall.Receives.State).To(Equal(state))
		})

		Context("when config does not contain system domain", func() {
			It("passes system domain from the state", func() {
				config := commands.GCPCreateLBsConfig{
					CertPath: "some-cert-path",
					KeyPath:  "some-key-path",
					LBType:   "cf",
					Domain:   "",
				}
				err := command.Execute(config, state)

				Expect(err).NotTo(HaveOccurred())
				Expect(gcpCreateLBs.ExecuteCall.CallCount).To(Equal(1))
				Expect(gcpCreateLBs.ExecuteCall.Receives.Config).To(Equal(commands.GCPCreateLBsConfig{
					CertPath: "some-cert-path",
					KeyPath:  "some-key-path",
					LBType:   "cf",
					Domain:   "some-domain",
				}))
				Expect(gcpCreateLBs.ExecuteCall.Receives.State).To(Equal(state))
			})
		})
	})
})
