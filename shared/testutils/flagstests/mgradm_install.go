// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flagstests

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/testutils"
)

// InstallFlagsTestArgs is the slice of command parameters to use with AssertInstallFlags.
var InstallFlagsTestArgs = func() []string {
	args := []string{
		"--tz", "CEST",
		"--email", "admin@foo.bar",
		"--emailfrom", "sender@foo.bar",
		"--issParent", "parent.iss.com",
		"--tftp=false",
		"--reportdb-user", "reportdbuser",
		"--reportdb-password", "reportdbpass",
		"--reportdb-name", "reportdbname",
		"--reportdb-host", "reportdbhost",
		"--reportdb-port", "5678",
		"--debug-java",
		"--admin-login", "adminuser",
		"--admin-password", "adminpass",
		"--admin-firstName", "adminfirst",
		"--admin-lastName", "adminlast",
		"--organization", "someorg",
	}

	args = append(args, SCCFlagTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)
	args = append(args, SalineFlagsTestArgs...)
	args = append(args, SSLGenerationFlagsTestArgs...)
	args = append(args, PgsqlFlagsTestArgs...)
	args = append(args, DBFlagsTestArgs...)
	args = append(args, ReportDBFlagsTestArgs...)
	args = append(args, InstallSSLFlagsTestArgs...)

	return args
}

// AssertInstallFlags checks that all the install flags are parsed correctly.
func AssertInstallFlags(t *testing.T, flags *utils.ServerFlags) {
	testutils.AssertEquals(t, "Error parsing --tz", "CEST", flags.Installation.TZ)
	testutils.AssertEquals(t, "Error parsing --email", "admin@foo.bar", flags.Installation.Email)
	testutils.AssertEquals(t, "Error parsing --emailfrom", "sender@foo.bar", flags.Installation.EmailFrom)
	testutils.AssertEquals(t, "Error parsing --issParent", "parent.iss.com", flags.Installation.IssParent)
	testutils.AssertEquals(t, "Error parsing --tftp", false, flags.Installation.Tftp)
	testutils.AssertTrue(t, "Error parsing --debug-java", flags.Installation.Debug.Java)
	testutils.AssertEquals(t, "Error parsing --admin-login", "adminuser", flags.Installation.Admin.Login)
	testutils.AssertEquals(t, "Error parsing --admin-password", "adminpass", flags.Installation.Admin.Password)
	testutils.AssertEquals(t, "Error parsing --admin-firstName", "adminfirst", flags.Installation.Admin.FirstName)
	testutils.AssertEquals(t, "Error parsing --admin-lastName", "adminlast", flags.Installation.Admin.LastName)
	testutils.AssertEquals(t, "Error parsing --organization", "someorg", flags.Installation.Organization)
	AssertSCCFlag(t, &flags.Installation.SCC)
	AssertImageFlag(t, &flags.Image)
	AssertCocoFlag(t, &flags.Coco)
	AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
	AssertSalineFlag(t, &flags.Saline)
	AssertPgsqlFlag(t, &flags.Pgsql)
	AssertDBFlag(t, &flags.Installation.DB)
	AssertReportDBFlag(t, &flags.Installation.ReportDB)
	AssertInstallSSLFlag(t, &flags.Installation.SSL)
}
