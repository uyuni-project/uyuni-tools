// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

// MigrateFlagsTestArgs is the slice of command parameters to use with AssertMigrateFlags.
var MigrateFlagsTestArgs = func() []string {
	args := []string{
		"source.fq.dn",
		"--mirror", "/path/to/mirror",
		"--tz", "CEST",
		"--prepare",
		"--user", "sudoer",
	}

	args = append(args, MirrorFlagTestArgs...)
	args = append(args, SCCFlagTestArgs...)
	args = append(args, PgsqlFlagsTestArgs...)
	args = append(args, DBFlagsTestArgs...)
	args = append(args, ReportDBFlagsTestArgs...)
	args = append(args, SSLGenerationFlagsTestArgs...)
	args = append(args, SalineFlagsTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, DBUpgradeImageFlagTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)
	args = append(args, InstallDBSSLFlagsTestArgs...)
	return args
}

// AssertMigrateFlags checks that all the migrate flags are parsed correctly.
func AssertMigrateFlags(t *testing.T, flags *utils.MigrationFlags) {
	testutils.AssertTrue(t, "Prepare not set", flags.Prepare)
	testutils.AssertEquals(t, "Error parsing --user", "sudoer", flags.User)
}
