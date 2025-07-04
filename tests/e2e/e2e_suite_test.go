// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/uyuni-project/uyuni-tools/tests/e2e/utils"
)

func TestE2E(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "End-To-End Test Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	pwd, err := os.Getwd()
	gomega.Expect(err).ToNot(gomega.HaveOccurred())

	// Initialize the exported variables from the utils package
	utils.MgradmPath = filepath.Join(pwd, "..", "..", "bin", "mgradm")
	utils.MgrctlPath = filepath.Join(pwd, "..", "..", "bin", "mgrctl")
	utils.MgradmConfigPath = filepath.Join(pwd, "utils", "mgradm.yaml")

	_, err = os.Stat(utils.MgradmPath)
	gomega.Expect(err).To(gomega.Succeed(), "mgradm binary not found at %s. Error: %v", utils.MgradmPath, err)

	_, err = os.Stat(utils.MgrctlPath)
	gomega.Expect(err).To(gomega.Succeed(), "mgrctl binary not found at %s. Error: %v", utils.MgrctlPath, err)

	ginkgo.GinkgoWriter.Printf("Using mgradm: %s\n", utils.MgradmPath)
	ginkgo.GinkgoWriter.Printf("Using mgrctl: %s\n", utils.MgrctlPath)
})
