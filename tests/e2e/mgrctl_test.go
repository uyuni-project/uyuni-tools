// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/uyuni-project/uyuni-tools/tests/e2e/utils"
)

var _ = ginkgo.Describe("mgrctl tests", func() {
	ginkgo.Context("Help command", func() {
		ginkgo.It("should show help for mgrctl command", func() {
			args := []string{"help"}
			output, err := utils.RunMgrctlCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(output).To(gomega.ContainSubstring("Tool to help managing Uyuni servers mainly through their API"))
		})
	})
	ginkgo.Context("API command", func() {
		ginkgo.It("should show help for term command", func() {
			args := []string{"api", "--help"}
			_, err := utils.RunMgrctlCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
	ginkgo.Context("Copy Command", func() {
		ginkgo.It("should show help for term command", func() {
			args := []string{"cp", "--help"}
			_, err := utils.RunMgrctlCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
	ginkgo.Context("Exec Command", func() {
		ginkgo.It("should show help for term command", func() {
			args := []string{"exec", "--help"}
			_, err := utils.RunMgrctlCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
	ginkgo.Context("Proxy Command", func() {
		ginkgo.It("should show help for term command", func() {
			args := []string{"proxy", "--help"}
			_, err := utils.RunMgrctlCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
	ginkgo.Context("Term Command", func() {
		ginkgo.It("should show help for term command", func() {
			args := []string{"term", "--help"}
			_, err := utils.RunMgrctlCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
})
