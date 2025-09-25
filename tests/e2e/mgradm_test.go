// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/uyuni-project/uyuni-tools/tests/e2e/utils"
)

const backend = "podman"

var _ = ginkgo.Describe("mgradm tests", func() {
	// Ensure that the environment is clean before each test.
	ginkgo.BeforeEach(func() {
		ginkgo.By("Running mgradm cleanup: uninstall --force --purge-volumes")
		cleanupArgs := []string{"uninstall", "--force", "--purge-volumes"}
		_, _ = utils.RunMgradmCommand(cleanupArgs)
		// gomega.Expect(err).To(gomega.Succeed(), "Failed to run mgradm cleanup command. Error: %v\nOutput: %s", err, output)
	})

	ginkgo.Context("Installation", func() {
		ginkgo.It("should install using a configuration YAML file", func() {
			args := []string{"install", backend, "--config", utils.MgradmConfigPath, "--logLevel", "debug"}
			output, err := utils.RunMgradmCommand(args)
			gomega.Expect(err).To(gomega.Succeed(), "Command failed with error: %v\nOutput: %s", err, output)
			gomega.Expect(output).Should(gomega.ContainSubstring("Setting up uyuni network"))
			gomega.Expect(output).Should(gomega.ContainSubstring("Enabling system service"))
			gomega.Expect(output).Should(gomega.ContainSubstring("Waiting for the server to start"))
			gomega.Expect(output).Should(gomega.ContainSubstring("Run setup command in the container"))
			gomega.Expect(output).Should(gomega.ContainSubstring("Populating the database"))
		})

		ginkgo.It("should show help for install command", func() {
			args := []string{"install", "--help"}
			output, err := utils.RunMgradmCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(output).To(gomega.ContainSubstring("Install a new server"))
		})
	})
})
