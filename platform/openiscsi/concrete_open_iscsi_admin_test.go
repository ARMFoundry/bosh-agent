package openiscsi_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-agent/platform/openiscsi"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
)

var _ = Describe("concreteOpenIscsiAdmin", func() {
	var (
		fs       *fakesys.FakeFileSystem
		runner   *fakesys.FakeCmdRunner
		iscsiAdm OpenIscsi
	)

	BeforeEach(func() {
		runner = fakesys.NewFakeCmdRunner()
		fs = fakesys.NewFakeFileSystem()
		iscsiAdm = NewConcreteOpenIscsiAdmin(fs, runner, boshlog.NewLogger(boshlog.LevelNone))
	})

	Describe("Setup", func() {
		var (
			initiatorName string
			username      string
			target        string
			password      string
		)
		BeforeEach(func() {
			initiatorName = "iqn.2007-05.com.fake-domain:fake-username"
			username = "fake-username"
			target = "11.11.22.22"
			password = "fake-password"
		})

		It("Open-iscsi setup successfully", func() {
			err := iscsiAdm.Setup(initiatorName, username, password)
			Expect(err).NotTo(HaveOccurred())

			initiatorNameIscsi := fs.GetFileTestStat("/etc/iscsi/initiatorname.iscsi")
			Expect(initiatorNameIscsi).ToNot(BeNil())

			expectedInitiatorNameIscsi := "InitiatorName=iqn.2007-05.com.fake-domain:fake-username"

			Expect(initiatorNameIscsi.StringContents()).To(Equal(expectedInitiatorNameIscsi))

			iscsidConf := fs.GetFileTestStat("/etc/iscsi/iscsid.conf")
			Expect(iscsidConf).ToNot(BeNil())

			expectedIscsidConf := `# Generated by bosh-agent
node.startup = automatic
node.session.auth.authmethod = CHAP
node.session.auth.username = fake-username
node.session.auth.password = fake-password
discovery.sendtargets.auth.authmethod = CHAP
discovery.sendtargets.auth.username = fake-username
discovery.sendtargets.auth.password = fake-password
node.session.timeo.replacement_timeout = 120
node.conn[0].timeo.login_timeout = 15
node.conn[0].timeo.logout_timeout = 15
node.conn[0].timeo.noop_out_interval = 10
node.conn[0].timeo.noop_out_timeout = 15
node.session.iscsi.InitialR2T = No
node.session.iscsi.ImmediateData = Yes
node.session.iscsi.FirstBurstLength = 262144
node.session.iscsi.MaxBurstLength = 16776192
node.conn[0].iscsi.MaxRecvDataSegmentLength = 65536
`
			Expect(iscsidConf.StringContents()).To(Equal(expectedIscsidConf))

		})

		It("Returns error when write to /etc/iscsi/initiatorname.iscsi fails", func() {
			fs.WriteFileErrors = map[string]error{
				"/etc/iscsi/initiatorname.iscsi": errors.New("fake-fs-error"),
			}
			err := iscsiAdm.Setup(initiatorName, username, password)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Writing to /etc/iscsi/initiatorname.iscsi"))
		})

		It("Returns error when write to /etc/iscsi/iscsid.conf fails", func() {
			fs.WriteFileErrors = map[string]error{
				"/etc/iscsi/iscsid.conf": errors.New("fake-fs-error"),
			}
			err := iscsiAdm.Setup(initiatorName, username, password)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Writing to /etc/iscsi/iscsid.conf"))
		})

		It("Returns error when run 'multipath-tools' fails", func() {
			runner.AddCmdResult(
				"/etc/init.d/multipath-tools restart",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Setup(initiatorName, username, password)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Restarting multipath after restarting open-iscsi"))
		})
	})

	Describe("Start", func() {
		It("Open-iscsi start successfully", func() {
			err := iscsiAdm.Start()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when run 'open-iscsi start' fails", func() {
			runner.AddCmdResult(
				"/etc/init.d/open-iscsi start",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Start()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-cmd-error"))
		})
	})

	Describe("Stop", func() {
		It("Open-iscsi stop successfully", func() {
			err := iscsiAdm.Stop()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when run 'open-iscsi stop' fails", func() {
			runner.AddCmdResult(
				"/etc/init.d/open-iscsi stop",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Stop()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-cmd-error"))
		})
	})

	Describe("Restart", func() {
		It("Open-iscsi restart successfully", func() {
			err := iscsiAdm.Restart()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when run 'open-iscsi restart' fails", func() {
			runner.AddCmdResult(
				"/etc/init.d/open-iscsi restart",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Restart()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-cmd-error"))
		})
	})

	Describe("Discovery", func() {
		var (
			target string
		)
		BeforeEach(func() {
			target = "11.11.22.22"
		})

		It("iscsiadm stop successfully", func() {
			err := iscsiAdm.Discovery(target)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when run 'iscsiadm discovery' fails", func() {
			runner.AddCmdResult(
				"iscsiadm -m discovery -t sendtargets -p 11.11.22.22",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Discovery(target)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-cmd-error"))
		})
	})

	Describe("Login", func() {
		It("iscsiadm login successfully", func() {
			err := iscsiAdm.Login()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when run 'iscsiadm login' fails", func() {
			runner.AddCmdResult(
				"iscsiadm -m node -l",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Login()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-cmd-error"))
		})
	})

	Describe("Logout", func() {
		It("iscsiadm logout successfully", func() {
			err := iscsiAdm.Logout()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns error when run 'iscsiadm logout' fails", func() {
			runner.AddCmdResult(
				"iscsiadm -m node -u",
				fakesys.FakeCmdResult{Error: errors.New("fake-cmd-error")},
			)
			err := iscsiAdm.Logout()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-cmd-error"))
		})
	})
})