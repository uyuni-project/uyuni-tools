// SPDX-FileCopyrightText: 2024 SUSE LLC
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
		"--db-user", "dbuser",
		"--db-password", "dbpass",
		"--db-name", "dbname",
		"--db-host", "dbhost",
		"--db-port", "1234",
		"--db-admin-user", "dbadmin",
		"--db-admin-password", "dbadminpass",
		"--db-provider", "aws",
		"--tftp=false",
		"--reportdb-user", "reportdbuser",
		"--reportdb-password", "reportdbpass",
		"--reportdb-name", "reportdbname",
		"--reportdb-host", "reportdbhost",
		"--reportdb-port", "5678",
		"--ssl-password", "sslsecret",
		"--ssl-ca-intermediate", "path/inter1.crt",
		"--ssl-ca-intermediate", "path/inter2.crt",
		"--ssl-ca-root", "path/root.crt",
		"--ssl-server-cert", "path/srv.crt",
		"--ssl-server-key", "path/srv.key",
		"--debug-java",
		"--admin-login", "adminuser",
		"--admin-password", "adminpass",
		"--admin-firstName", "adminfirst",
		"--admin-lastName", "adminlast",
		"--organization", "someorg",
	}

	args = append(args, MirrorFlagTestArgs...)
	args = append(args, SCCFlagTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)
	args = append(args, SalineFlagsTestArgs...)
	args = append(args, SSLGenerationFlagsTestArgs...)

	return args
}

// AssertInstallFlags checks that all the install flags are parsed correctly.
func AssertInstallFlags(t *testing.T, flags *utils.ServerFlags) {
	testutils.AssertEquals(t, "Error parsing --tz", "CEST", flags.Installation.TZ)
	testutils.AssertEquals(t, "Error parsing --email", "admin@foo.bar", flags.Installation.Email)
	testutils.AssertEquals(t, "Error parsing --emailfrom", "sender@foo.bar", flags.Installation.EmailFrom)
	testutils.AssertEquals(t, "Error parsing --issParent", "parent.iss.com", flags.Installation.IssParent)
	testutils.AssertEquals(t, "Error parsing --db-user", "dbuser", flags.Installation.DB.User)
	testutils.AssertEquals(t, "Error parsing --db-password", "dbpass", flags.Installation.DB.Password)
	testutils.AssertEquals(t, "Error parsing --db-name", "dbname", flags.Installation.DB.Name)
	testutils.AssertEquals(t, "Error parsing --db-host", "dbhost", flags.Installation.DB.Host)
	testutils.AssertEquals(t, "Error parsing --db-port", 1234, flags.Installation.DB.Port)
	testutils.AssertEquals(t, "Error parsing --db-admin-user", "dbadmin", flags.Installation.DB.Admin.User)
	testutils.AssertEquals(t, "Error parsing --db-admin-password", "dbadminpass", flags.Installation.DB.Admin.Password)
	testutils.AssertEquals(t, "Error parsing --db-provider", "aws", flags.Installation.DB.Provider)
	testutils.AssertEquals(t, "Error parsing --tftp", false, flags.Installation.Tftp)
	testutils.AssertEquals(t, "Error parsing --reportdb-user", "reportdbuser", flags.Installation.ReportDB.User)
	testutils.AssertEquals(t, "Error parsing --reportdb-password", "reportdbpass", flags.Installation.ReportDB.Password)
	testutils.AssertEquals(t, "Error parsing --reportdb-name", "reportdbname", flags.Installation.ReportDB.Name)
	testutils.AssertEquals(t, "Error parsing --reportdb-host", "reportdbhost", flags.Installation.ReportDB.Host)
	testutils.AssertEquals(t, "Error parsing --reportdb-port", 5678, flags.Installation.ReportDB.Port)
	testutils.AssertEquals(t, "Error parsing --ssl-cname", []string{"cname1", "cname2"}, flags.Installation.SSL.Cnames)
	testutils.AssertEquals(t, "Error parsing --ssl-country", "OS", flags.Installation.SSL.Country)
	testutils.AssertEquals(t, "Error parsing --ssl-state", "sslstate", flags.Installation.SSL.State)
	testutils.AssertEquals(t, "Error parsing --ssl-city", "sslcity", flags.Installation.SSL.City)
	testutils.AssertEquals(t, "Error parsing --ssl-org", "sslorg", flags.Installation.SSL.Org)
	testutils.AssertEquals(t, "Error parsing --ssl-ou", "sslou", flags.Installation.SSL.OU)
	testutils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", flags.Installation.SSL.Password)
	testutils.AssertEquals(t, "Error parsing --ssl-ca-intermediate",
		[]string{"path/inter1.crt", "path/inter2.crt"}, flags.Installation.SSL.Ca.Intermediate,
	)
	testutils.AssertEquals(t, "Error parsing --ssl-ca-root", "path/root.crt", flags.Installation.SSL.Ca.Root)
	testutils.AssertEquals(t, "Error parsing --ssl-server-cert", "path/srv.crt", flags.Installation.SSL.Server.Cert)
	testutils.AssertEquals(t, "Error parsing --ssl-server-key", "path/srv.key", flags.Installation.SSL.Server.Key)
	testutils.AssertTrue(t, "Error parsing --debug-java", flags.Installation.Debug.Java)
	testutils.AssertEquals(t, "Error parsing --admin-login", "adminuser", flags.Installation.Admin.Login)
	testutils.AssertEquals(t, "Error parsing --admin-password", "adminpass", flags.Installation.Admin.Password)
	testutils.AssertEquals(t, "Error parsing --admin-firstName", "adminfirst", flags.Installation.Admin.FirstName)
	testutils.AssertEquals(t, "Error parsing --admin-lastName", "adminlast", flags.Installation.Admin.LastName)
	testutils.AssertEquals(t, "Error parsing --organization", "someorg", flags.Installation.Organization)
	AssertMirrorFlag(t, flags.Mirror)
	AssertSCCFlag(t, &flags.Installation.SCC)
	AssertImageFlag(t, &flags.Image)
	AssertCocoFlag(t, &flags.Coco)
	AssertHubXmlrpcFlag(t, &flags.HubXmlrpc)
	AssertSalineFlag(t, &flags.Saline)
}
