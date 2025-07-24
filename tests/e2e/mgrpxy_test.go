// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package e2e_test

import (
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/uyuni-project/uyuni-tools/tests/e2e/utils"
)

var _ = ginkgo.Describe("mgrpxy tests", func() {
	ginkgo.Context("Help command", func() {
		ginkgo.It("should show help for mgrpxy command", func() {
			args := []string{"help"}
			_, err := utils.RunMgrpxyCommand(args)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})
})
