// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// DbFlags can store all values required to connect to a database.
type DbFlags struct {
	Host     string
	Name     string
	Port     int
	User     string
	Password string
	Protocol string
	Provider string
	Admin    struct {
		User     string
		Password string
	}
}

// DebugFlags contains information about enabled/disabled debug.
type DebugFlags struct {
	Java bool
}

// InstallFlags stores all the flags used by install command.
type InstallFlags struct {
	TZ           string
	Email        string
	EmailFrom    string
	IssParent    string
	Mirror       string
	Tftp         bool
	Db           DbFlags
	ReportDb     DbFlags
	Ssl          cmd_utils.SslCertFlags
	Scc          types.SCCCredentials
	Debug        DebugFlags
	Image        types.ImageFlags `mapstructure:",squash"`
	Coco         cmd_utils.CocoFlags
	HubXmlrpc    cmd_utils.HubXmlrpcFlags
	Admin        apiTypes.User
	Organization string
}

// idChecker verifies that the value is a valid identifier.
func idChecker(value string) bool {
	r := regexp.MustCompile(`^([[:alnum:]]|[._-])+$`)
	if r.MatchString(value) {
		return true
	}
	fmt.Println(L("Can only contain letters, digits . _ and -"))
	return false
}

// emailChecker verifies that the value is a valid email address.
func emailChecker(value string) bool {
	address, err := mail.ParseAddress(value)
	if err != nil || address.Name != "" || strings.ContainsAny(value, "<>") {
		fmt.Println(L("Not a valid email address"))
		return false
	}
	return true
}

// CheckParameters checks parameters for install command.
func (flags *InstallFlags) CheckParameters(cmd *cobra.Command, command string) {
	if flags.Db.Password == "" {
		flags.Db.Password = utils.GetRandomBase64(30)
	}

	if flags.ReportDb.Password == "" {
		flags.ReportDb.Password = utils.GetRandomBase64(30)
	}

	// Make sure we have all the required 3rd party flags or none
	flags.Ssl.CheckParameters()

	// Since we use cert-manager for self-signed certificates on kubernetes we don't need password for it
	if !flags.Ssl.UseExisting() && command == "podman" {
		utils.AskPasswordIfMissing(&flags.Ssl.Password, cmd.Flag("ssl-password").Usage, 0, 0)
	}

	// Use the host timezone if the user didn't define one
	if flags.TZ == "" {
		flags.TZ = utils.GetLocalTimezone()
	}

	utils.AskIfMissing(&flags.Email, cmd.Flag("email").Usage, 1, 128, emailChecker)
	utils.AskIfMissing(&flags.EmailFrom, cmd.Flag("emailfrom").Usage, 0, 0, emailChecker)

	utils.AskIfMissing(&flags.Admin.Login, cmd.Flag("admin-login").Usage, 1, 64, idChecker)
	utils.AskPasswordIfMissing(&flags.Admin.Password, cmd.Flag("admin-password").Usage, 5, 48)
	utils.AskIfMissing(&flags.Organization, cmd.Flag("organization").Usage, 3, 128, nil)

	flags.Ssl.Email = flags.Email
	flags.Admin.Email = flags.Email
}

// AddInspectFlags add flags to inspect command.
func AddInspectFlags(cmd *cobra.Command) {
	cmd_utils.AddSCCFlag(cmd)
	cmd_utils.AddImageFlag(cmd)
}

// AddInstallFlags add flags to installa command.
func AddInstallFlags(cmd *cobra.Command) {
	cmd_utils.AddMirrorFlag(cmd)
	cmd.Flags().String("tz", "", L("Time zone to set on the server. Defaults to the host timezone"))
	cmd.Flags().String("email", "admin@example.com", L("Administrator e-mail"))
	cmd.Flags().String("emailfrom", "notifications@example.com", L("E-Mail sending the notifications"))
	cmd.Flags().String("issParent", "", L("InterServerSync v1 parent FQDN"))

	cmd.Flags().String("db-user", "spacewalk", L("Database user"))
	cmd.Flags().String("db-password", "", L("Database password. Randomly generated by default"))
	cmd.Flags().String("db-name", "susemanager", L("Database name"))
	cmd.Flags().String("db-host", "localhost", L("Database host"))
	cmd.Flags().Int("db-port", 5432, L("Database port"))
	cmd.Flags().String("db-protocol", "tcp", L("Database protocol"))
	cmd.Flags().String("db-admin-user", "", L("External database admin user name"))
	cmd.Flags().String("db-admin-password", "", L("External database admin password"))
	cmd.Flags().String("db-provider", "", L("External database provider. Possible values 'aws'"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "db", Title: L("Database Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "db-user", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-password", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-name", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-host", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-port", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-protocol", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-admin-user", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-admin-password", "db")
	_ = utils.AddFlagToHelpGroupID(cmd, "db-provider", "db")

	cmd.Flags().Bool("tftp", true, L("Enable TFTP"))
	cmd.Flags().String("reportdb-name", "reportdb", L("Report database name"))
	cmd.Flags().String("reportdb-host", "localhost", L("Report database host"))
	cmd.Flags().Int("reportdb-port", 5432, L("Report database port"))
	cmd.Flags().String("reportdb-user", "pythia_susemanager", L("Report Database username"))
	cmd.Flags().String("reportdb-password", "", L("Report database password. Randomly generated by default"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "reportdb", Title: L("Report DB Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-name", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-host", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-port", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-user", "reportdb")
	_ = utils.AddFlagToHelpGroupID(cmd, "reportdb-password", "reportdb")

	// For generated CA and certificate
	cmd.Flags().StringSlice("ssl-cname", []string{}, L("SSL certificate cnames separated by commas"))
	cmd.Flags().String("ssl-country", "DE", L("SSL certificate country"))
	cmd.Flags().String("ssl-state", "Bayern", L("SSL certificate state"))
	cmd.Flags().String("ssl-city", "Nuernberg", L("SSL certificate city"))
	cmd.Flags().String("ssl-org", "SUSE", L("SSL certificate organization"))
	cmd.Flags().String("ssl-ou", "SUSE", L("SSL certificate organization unit"))
	cmd.Flags().String("ssl-password", "", L("Password for the CA key to generate"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "ssl", Title: L("SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-cname", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-country", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-state", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-city", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-org", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ou", "ssl")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-password", "ssl")

	// For SSL 3rd party certificates
	cmd.Flags().StringSlice("ssl-ca-intermediate", []string{}, L("Intermediate CA certificate path"))
	cmd.Flags().String("ssl-ca-root", "", L("Root CA certificate path"))
	cmd.Flags().String("ssl-server-cert", "", L("Server certificate path"))
	cmd.Flags().String("ssl-server-key", "", L("Server key path"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "ssl3rd", Title: L("3rd Party SSL Certificate Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-intermediate", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-ca-root", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-cert", "ssl3rd")
	_ = utils.AddFlagToHelpGroupID(cmd, "ssl-server-key", "ssl3rd")

	cmd_utils.AddSCCFlag(cmd)

	cmd.Flags().Bool("debug-java", false, L("Enable tomcat and taskomatic remote debugging"))
	cmd_utils.AddImageFlag(cmd)

	cmd_utils.AddCocoFlag(cmd)

	cmd_utils.AddHubXmlrpcFlags(cmd)

	cmd.Flags().String("admin-login", "admin", L("Administrator user name"))
	cmd.Flags().String("admin-password", "", L("Administrator password"))
	cmd.Flags().String("admin-firstName", "Administrator", L("First name of the administrator"))
	cmd.Flags().String("admin-lastName", "McAdmin", L("Last name of the administrator"))
	cmd.Flags().String("organization", "Organization", L("First organization name"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "first-user", Title: L("First User Flags")})
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-login", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-password", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-firstName", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "admin-lastName", "first-user")
	_ = utils.AddFlagToHelpGroupID(cmd, "organization", "first-user")
}
