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
		"srv.fq.dn",
		"--mirror", "/path/to/mirror",
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
	args = append(args, DBUpgradeImageFlagTestArgs...)
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
func AssertInstallFlags(t *testing.T, flags *utils.InstallationFlags) {
	testutils.AssertEquals(t, "Error parsing --email", "admin@foo.bar", flags.Email)
	testutils.AssertEquals(t, "Error parsing --emailfrom", "sender@foo.bar", flags.EmailFrom)
	testutils.AssertEquals(t, "Error parsing --issParent", "parent.iss.com", flags.IssParent)
	testutils.AssertEquals(t, "Error parsing --tftp", false, flags.Tftp)
	testutils.AssertTrue(t, "Error parsing --debug-java", flags.Debug.Java)
	testutils.AssertEquals(t, "Error parsing --admin-login", "adminuser", flags.Admin.Login)
	testutils.AssertEquals(t, "Error parsing --admin-password", "adminpass", flags.Admin.Password)
	testutils.AssertEquals(t, "Error parsing --admin-firstName", "adminfirst", flags.Admin.FirstName)
	testutils.AssertEquals(t, "Error parsing --admin-lastName", "adminlast", flags.Admin.LastName)
	testutils.AssertEquals(t, "Error parsing --organization", "someorg", flags.Organization)
	AssertSCCFlag(t, &flags.SCC)
	AssertDBFlag(t, &flags.DB)
	AssertReportDBFlag(t, &flags.ReportDB)
}
