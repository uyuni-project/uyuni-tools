// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/shared/test_utils"
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
		"--db-protocol", "dbprot",
		"--db-admin-user", "dbadmin",
		"--db-admin-password", "dbadminpass",
		"--db-provider", "aws",
		"--tftp=false",
		"--reportdb-user", "reportdbuser",
		"--reportdb-password", "reportdbpass",
		"--reportdb-name", "reportdbname",
		"--reportdb-host", "reportdbhost",
		"--reportdb-port", "5678",
		"--ssl-cname", "cname1",
		"--ssl-cname", "cname2",
		"--ssl-country", "OS",
		"--ssl-state", "sslstate",
		"--ssl-city", "sslcity",
		"--ssl-org", "sslorg",
		"--ssl-ou", "sslou",
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
	args = append(args, SccFlagTestArgs...)
	args = append(args, ImageFlagsTestArgs...)
	args = append(args, CocoFlagsTestArgs...)
	args = append(args, HubXmlrpcFlagsTestArgs...)

	return args
}

// AssertInstallFlags checks that all the install flags are parsed correctly.
func AssertInstallFlags(t *testing.T, cmd *cobra.Command, flags *shared.InstallFlags) {
	test_utils.AssertEquals(t, "Error parsing --tz", "CEST", flags.TZ)
	test_utils.AssertEquals(t, "Error parsing --email", "admin@foo.bar", flags.Email)
	test_utils.AssertEquals(t, "Error parsing --emailfrom", "sender@foo.bar", flags.EmailFrom)
	test_utils.AssertEquals(t, "Error parsing --issParent", "parent.iss.com", flags.IssParent)
	test_utils.AssertEquals(t, "Error parsing --db-user", "dbuser", flags.Db.User)
	test_utils.AssertEquals(t, "Error parsing --db-password", "dbpass", flags.Db.Password)
	test_utils.AssertEquals(t, "Error parsing --db-name", "dbname", flags.Db.Name)
	test_utils.AssertEquals(t, "Error parsing --db-host", "dbhost", flags.Db.Host)
	test_utils.AssertEquals(t, "Error parsing --db-port", 1234, flags.Db.Port)
	test_utils.AssertEquals(t, "Error parsing --db-protocol", "dbprot", flags.Db.Protocol)
	test_utils.AssertEquals(t, "Error parsing --db-admin-user", "dbadmin", flags.Db.Admin.User)
	test_utils.AssertEquals(t, "Error parsing --db-admin-password", "dbadminpass", flags.Db.Admin.Password)
	test_utils.AssertEquals(t, "Error parsing --db-provider", "aws", flags.Db.Provider)
	test_utils.AssertEquals(t, "Error parsing --tftp", false, flags.Tftp)
	test_utils.AssertEquals(t, "Error parsing --reportdb-user", "reportdbuser", flags.ReportDb.User)
	test_utils.AssertEquals(t, "Error parsing --reportdb-password", "reportdbpass", flags.ReportDb.Password)
	test_utils.AssertEquals(t, "Error parsing --reportdb-name", "reportdbname", flags.ReportDb.Name)
	test_utils.AssertEquals(t, "Error parsing --reportdb-host", "reportdbhost", flags.ReportDb.Host)
	test_utils.AssertEquals(t, "Error parsing --reportdb-port", 5678, flags.ReportDb.Port)
	test_utils.AssertEquals(t, "Error parsing --ssl-cname", []string{"cname1", "cname2"}, flags.Ssl.Cnames)
	test_utils.AssertEquals(t, "Error parsing --ssl-country", "OS", flags.Ssl.Country)
	test_utils.AssertEquals(t, "Error parsing --ssl-state", "sslstate", flags.Ssl.State)
	test_utils.AssertEquals(t, "Error parsing --ssl-city", "sslcity", flags.Ssl.City)
	test_utils.AssertEquals(t, "Error parsing --ssl-org", "sslorg", flags.Ssl.Org)
	test_utils.AssertEquals(t, "Error parsing --ssl-ou", "sslou", flags.Ssl.OU)
	test_utils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", flags.Ssl.Password)
	test_utils.AssertEquals(t, "Error parsing --ssl-ca-intermediate",
		[]string{"path/inter1.crt", "path/inter2.crt"}, flags.Ssl.Ca.Intermediate,
	)
	test_utils.AssertEquals(t, "Error parsing --ssl-ca-root", "path/root.crt", flags.Ssl.Ca.Root)
	test_utils.AssertEquals(t, "Error parsing --ssl-server-cert", "path/srv.crt", flags.Ssl.Server.Cert)
	test_utils.AssertEquals(t, "Error parsing --ssl-server-key", "path/srv.key", flags.Ssl.Server.Key)
	test_utils.AssertTrue(t, "Error parsing --debug-java", flags.Debug.Java)
	test_utils.AssertEquals(t, "Error parsing --admin-login", "adminuser", flags.Admin.Login)
	test_utils.AssertEquals(t, "Error parsing --admin-password", "adminpass", flags.Admin.Password)
	test_utils.AssertEquals(t, "Error parsing --admin-firstName", "adminfirst", flags.Admin.FirstName)
	test_utils.AssertEquals(t, "Error parsing --admin-lastName", "adminlast", flags.Admin.LastName)
	test_utils.AssertEquals(t, "Error parsing --organization", "someorg", flags.Organization)
	AssertMirrorFlag(t, cmd, flags.Mirror)
	AssertSccFlag(t, cmd, &flags.Scc)
	AssertImageFlag(t, cmd, &flags.Image)
	AssertCocoFlag(t, cmd, &flags.Coco)
	AssertHubXmlrpcFlag(t, cmd, &flags.HubXmlrpc)
}
