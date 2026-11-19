// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
)

// UpgradeFlagsTestArgs is the slice of command parameters to use with AssertUpgradeFlags.
var UpgradeFlagsTestArgs = func() []string {
	args := []string{}
	args = append(args, SCCFlagTestArgs...)
	args = append(args, PgsqlFlagsTestArgs...)
	args = append(args, SalineFlagsTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, RegistryImageFlagsTestArgs...)
	args = append(args, DBUpdateImageFlagTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)
	args = append(args, TFTPDFlagsTestArgs...)
	return args
}

// AssertUpgradeFlags checks that all the upgrade flags are parsed correctly.
func AssertUpgradeFlags(t *testing.T, flags *utils.ServerFlags) {
	AssertImageFlag(t, &flags.Image)
	AssertRegistryFlag(t, &flags.Image.Registry)
	AssertDBUpgradeImageFlag(t, &flags.DBUpgradeImage)
	AssertCocoFlag(t, &flags.Coco)
	AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
	AssertSalineFlag(t, &flags.Saline)
	AssertSCCFlag(t, &flags.Installation.SCC)
	AssertPgsqlFlag(t, &flags.Pgsql)
	AssertTFTPDFlag(t, &flags.TFTPD)
}
