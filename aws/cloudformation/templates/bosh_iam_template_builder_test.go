package templates_test

import (
	"github.com/cloudfoundry/bosh-bootloader/aws/cloudformation/templates"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BOSHIAMTemplateBuilder", func() {
	var (
		builder templates.BOSHIAMTemplateBuilder
	)

	BeforeEach(func() {
		builder = templates.NewBOSHIAMTemplateBuilder()
	})

	Describe("BOSHIAMUser", func() {
		Context("when we create a new bbl that supports tagging", func() {
			It("returns a template with Username", func() {
				user := builder.BOSHIAMUser("bosh-iam-user-bbl-env-lake-name-2016-08-15-12-03-00")
				Expect(user.Resources).To(HaveLen(2))
				IAMUser := user.Resources["BOSHUser"].Properties.(templates.IAMUser)
				Expect(IAMUser.UserName).To(Equal("bosh-iam-user-bbl-env-lake-name-2016-08-15-12-03-00"))
			})
		})

		Context("when we create a new bbl that does not support tagging", func() {
			It("returns a template for a BOSH IAM user", func() {
				user := builder.BOSHIAMUser("some-user-name")
				Expect(user.Resources).To(HaveLen(2))
				Expect(user.Resources).To(HaveKeyWithValue("BOSHUser", templates.Resource{
					Type: "AWS::IAM::User",
					Properties: templates.IAMUser{
						Policies: []templates.IAMPolicy{
							{
								PolicyName: "aws-cpi",
								PolicyDocument: templates.IAMPolicyDocument{
									Version: "2012-10-17",
									Statement: []templates.IAMStatement{
										{
											Action: []string{
												"ec2:AssociateAddress",
												"ec2:AttachVolume",
												"ec2:CreateVolume",
												"ec2:DeleteSnapshot",
												"ec2:DeleteVolume",
												"ec2:DescribeAddresses",
												"ec2:DescribeImages",
												"ec2:DescribeInstances",
												"ec2:DescribeRegions",
												"ec2:DescribeSecurityGroups",
												"ec2:DescribeSnapshots",
												"ec2:DescribeSubnets",
												"ec2:DescribeVolumes",
												"ec2:DetachVolume",
												"ec2:CreateSnapshot",
												"ec2:CreateTags",
												"ec2:RunInstances",
												"ec2:TerminateInstances",
												"ec2:RegisterImage",
												"ec2:DeregisterImage",
											},
											Effect:   "Allow",
											Resource: "*",
										},
										{
											Action:   []string{"elasticloadbalancing:*"},
											Effect:   "Allow",
											Resource: "*",
										},
									},
								},
							},
						},
					},
				}))

				Expect(user.Resources).To(HaveKeyWithValue("BOSHUserAccessKey", templates.Resource{
					Properties: templates.IAMAccessKey{
						UserName: templates.Ref{"BOSHUser"},
					},
					Type: "AWS::IAM::AccessKey",
				}))

				Expect(user.Outputs).To(HaveLen(2))
				Expect(user.Outputs).To(HaveKeyWithValue("BOSHUserAccessKey", templates.Output{
					Value: templates.Ref{"BOSHUserAccessKey"},
				}))

				Expect(user.Outputs).To(HaveKeyWithValue("BOSHUserSecretAccessKey", templates.Output{
					Value: templates.FnGetAtt{
						[]string{
							"BOSHUserAccessKey",
							"SecretAccessKey",
						},
					},
				}))
			})
		})
	})
})
