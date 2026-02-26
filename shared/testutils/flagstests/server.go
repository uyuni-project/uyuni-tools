// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
)

// ServerFlagsTestArgs is the slide of server-related command parameters to use with AssertServerFlags.
var ServerFlagsTestArgs = func() []string {
	args := []string{}
	args = append(args, SCCFlagTestArgs...)
	args = append(args, PgsqlFlagsTestArgs...)
	args = append(args, DBFlagsTestArgs...)
	args = append(args, ReportDBFlagsTestArgs...)
	args = append(args, InstallSSLFlagsTestArgs...)
	args = append(args, InstallDBSSLFlagsTestArgs...)
	args = append(args, SSLGenerationFlagsTestArgs...)
	args = append(args, SalineFlagsTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, RegistryImageFlagsTestArgs...)
	args = append(args, DBUpdateImageFlagTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)
	args = append(args, TFTPDFlagsTestArgs...)
	return args
}

// AssertServerFlags checks that all the server-related common flags are parsed correctly.
func AssertServerFlags(t *testing.T, flags *utils.ServerFlags) {
	AssertImageFlag(t, &flags.Image)
	AssertRegistryFlag(t, &flags.Image.Registry)
	AssertDBUpgradeImageFlag(t, &flags.DBUpgradeImage)
	AssertCocoFlag(t, &flags.Coco)
	AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
	AssertSalineFlag(t, &flags.Saline)
	AssertSCCFlag(t, &flags.Installation.SCC)
	AssertPgsqlFlag(t, &flags.Pgsql)
	AssertDBFlag(t, &flags.Installation.DB)
	AssertReportDBFlag(t, &flags.Installation.ReportDB)
	AssertInstallDBSSLFlag(t, &flags.Installation.SSL.DB)
	AssertInstallSSLFlag(t, &flags.Installation.SSL)
	AssertSSLGenerationFlag(t, &flags.Installation.SSL.SSLCertGenerationFlags)
	AssertTFTPDFlag(t, &flags.TFTPD)
}
