// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package flags_tests

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
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
func AssertInstallFlags(t *testing.T, cmd *cobra.Command, flags *utils.ServerFlags) {
	installFlags := flags.Installation
	test_utils.AssertEquals(t, "Error parsing --tz", "CEST", installFlags.TZ)
	test_utils.AssertEquals(t, "Error parsing --email", "admin@foo.bar", installFlags.Email)
	test_utils.AssertEquals(t, "Error parsing --emailfrom", "sender@foo.bar", installFlags.EmailFrom)
	test_utils.AssertEquals(t, "Error parsing --issParent", "parent.iss.com", installFlags.IssParent)
	test_utils.AssertEquals(t, "Error parsing --db-user", "dbuser", installFlags.Db.User)
	test_utils.AssertEquals(t, "Error parsing --db-password", "dbpass", installFlags.Db.Password)
	test_utils.AssertEquals(t, "Error parsing --db-name", "dbname", installFlags.Db.Name)
	test_utils.AssertEquals(t, "Error parsing --db-host", "dbhost", installFlags.Db.Host)
	test_utils.AssertEquals(t, "Error parsing --db-port", 1234, installFlags.Db.Port)
	test_utils.AssertEquals(t, "Error parsing --db-protocol", "dbprot", installFlags.Db.Protocol)
	test_utils.AssertEquals(t, "Error parsing --db-admin-user", "dbadmin", installFlags.Db.Admin.User)
	test_utils.AssertEquals(t, "Error parsing --db-admin-password", "dbadminpass", installFlags.Db.Admin.Password)
	test_utils.AssertEquals(t, "Error parsing --db-provider", "aws", installFlags.Db.Provider)
	test_utils.AssertEquals(t, "Error parsing --tftp", false, installFlags.Tftp)
	test_utils.AssertEquals(t, "Error parsing --reportdb-user", "reportdbuser", installFlags.ReportDb.User)
	test_utils.AssertEquals(t, "Error parsing --reportdb-password", "reportdbpass", installFlags.ReportDb.Password)
	test_utils.AssertEquals(t, "Error parsing --reportdb-name", "reportdbname", installFlags.ReportDb.Name)
	test_utils.AssertEquals(t, "Error parsing --reportdb-host", "reportdbhost", installFlags.ReportDb.Host)
	test_utils.AssertEquals(t, "Error parsing --reportdb-port", 5678, installFlags.ReportDb.Port)
	test_utils.AssertEquals(t, "Error parsing --ssl-cname", []string{"cname1", "cname2"}, installFlags.Ssl.Cnames)
	test_utils.AssertEquals(t, "Error parsing --ssl-country", "OS", installFlags.Ssl.Country)
	test_utils.AssertEquals(t, "Error parsing --ssl-state", "sslstate", installFlags.Ssl.State)
	test_utils.AssertEquals(t, "Error parsing --ssl-city", "sslcity", installFlags.Ssl.City)
	test_utils.AssertEquals(t, "Error parsing --ssl-org", "sslorg", installFlags.Ssl.Org)
	test_utils.AssertEquals(t, "Error parsing --ssl-ou", "sslou", installFlags.Ssl.OU)
	test_utils.AssertEquals(t, "Error parsing --ssl-password", "sslsecret", installFlags.Ssl.Password)
	test_utils.AssertEquals(t, "Error parsing --ssl-ca-intermediate",
		[]string{"path/inter1.crt", "path/inter2.crt"}, installFlags.Ssl.Ca.Intermediate,
	)
	test_utils.AssertEquals(t, "Error parsing --ssl-ca-root", "path/root.crt", installFlags.Ssl.Ca.Root)
	test_utils.AssertEquals(t, "Error parsing --ssl-server-cert", "path/srv.crt", installFlags.Ssl.Server.Cert)
	test_utils.AssertEquals(t, "Error parsing --ssl-server-key", "path/srv.key", installFlags.Ssl.Server.Key)
	test_utils.AssertTrue(t, "Error parsing --debug-java", installFlags.Debug.Java)
	test_utils.AssertEquals(t, "Error parsing --admin-login", "adminuser", installFlags.Admin.Login)
	test_utils.AssertEquals(t, "Error parsing --admin-password", "adminpass", installFlags.Admin.Password)
	test_utils.AssertEquals(t, "Error parsing --admin-firstName", "adminfirst", installFlags.Admin.FirstName)
	test_utils.AssertEquals(t, "Error parsing --admin-lastName", "adminlast", installFlags.Admin.LastName)
	test_utils.AssertEquals(t, "Error parsing --organization", "someorg", installFlags.Organization)
	AssertMirrorFlag(t, cmd, flags.Mirror)
	AssertSccFlag(t, cmd, &installFlags.Scc)
	AssertImageFlag(t, cmd, &flags.Image)
	AssertCocoFlag(t, cmd, &flags.Coco)
	AssertHubXmlrpcFlag(t, cmd, &flags.HubXmlrpc)
}
