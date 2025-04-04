// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
)

// UpgradeFlagsTestArgs is the slice of command parameters to use with AssertUpgradeFlags.
var UpgradeFlagsTestArgs = func() []string {
	args := []string{
		"--mirror", "/path/to/mirror",
		"--tz", "CEST",
	}
	args = append(args, MirrorFlagTestArgs...)
	args = append(args, SCCFlagTestArgs...)
	args = append(args, PgsqlFlagsTestArgs...)
	args = append(args, DBFlagsTestArgs...)
	args = append(args, ReportDBFlagsTestArgs...)
	args = append(args, SSLGenerationFlagsTestArgs...)
	args = append(args, InstallDBSSLFlagsTestArgs...)
	args = append(args, SalineFlagsTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, DBUpgradeImageFlagTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)
	return args
}

// AssertInstallFlags checks that all the install flags are parsed correctly.
func AssertUpgradeFlags(t *testing.T, flags *utils.UpgradeFlags) {
	AssertSCCFlag(t, &flags.SCC)
	AssertDBFlag(t, &flags.DB)
	AssertReportDBFlag(t, &flags.ReportDB)
}
