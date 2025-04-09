// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ServerFlags is a structure hosting the parameters for installation, migration and upgrade.
type ServerFlags struct {
	Mirror         string
	TZ             string
	DB             DBFlags
	DBUpgradeImage types.ImageFlags `mapstructure:"dbupgrade"`
	Image          types.ImageFlags `mapstructure:",squash"`
	Coco           CocoFlags
	HubXmlrpc      HubXmlrpcFlags
	Saline         SalineFlags
	Pgsql          types.PgsqlFlags
}

// DBFlags can store all values required to connect to a database.
type DBFlags struct {
	Host     string
	Name     string
	Port     int
	User     string
	Password string
	Provider string
	Admin    struct {
		User     string
		Password string
	}
}

// InstallationFlags contains the parameters that are used only for the installation of a new server.
type InstallationFlags struct {
	Email        string
	EmailFrom    string
	IssParent    string
	Tftp         bool
	DB           DBFlags
	ReportDB     DBFlags
	SSL          *InstallSSLFlags `mapstructure:"ssl"`
	SCC          types.SCCCredentials
	Debug        DebugFlags
	Admin        apiTypes.User
	Organization string
}

// MigrationFlags contains the parameters that are used only for the migration of a new server.
type MigrationFlags struct {
	// Prepare defines whether to run the full migration or just the data synchronization.
	Prepare bool
	// SourceUser is the username to use to connect to the source server in a migration.
	User string
}

// UpgradeFlags contains the parameters that are used only for the upgrade of a new server.
type UpgradeFlags struct {
	IssParent string
	Tftp      bool
	DB        DBFlags
	ReportDB  DBFlags
	SSL       *UpgradeSSLFlags `mapstructure:"ssl"`
	SCC       types.SCCCredentials
	Debug     DebugFlags
}

// CheckParameters checks parameters for install command.
func (flags *InstallationFlags) CheckParameters(cmd *cobra.Command, command string) {
	flags.SetPasswordIfMissing()

	flags.CheckSSLParameters(cmd, command)
	utils.AskIfMissing(&flags.Email, cmd.Flag("email").Usage, 1, 128, emailChecker)
	utils.AskIfMissing(&flags.EmailFrom, cmd.Flag("emailfrom").Usage, 0, 0, emailChecker)

	utils.AskIfMissing(&flags.Admin.Login, cmd.Flag("admin-login").Usage, 1, 64, idChecker)
	utils.AskPasswordIfMissing(&flags.Admin.Password, cmd.Flag("admin-password").Usage, 5, 48)
	utils.AskIfMissing(&flags.Organization, cmd.Flag("organization").Usage, 3, 128, nil)

	flags.SSL.Email = flags.Email
	flags.Admin.Email = flags.Email
}

func (flags *UpgradeFlags) CheckParameters(cmd *cobra.Command, command string) {
	flags.SetPasswordIfMissing()

	flags.CheckSSLParameters(cmd, command)
}

func (flags *DBFlags) SetPasswordIfMissing() {
	if flags.Password == "" {
		flags.Password = utils.GetRandomBase64(30)
	}
}

func (flags *InstallationFlags) SetPasswordIfMissing() {
	flags.DB.SetPasswordIfMissing()
	flags.ReportDB.SetPasswordIfMissing()

	// The admin password is only needed for local database
	if flags.DB.IsLocal() && flags.DB.Admin.Password == "" {
		flags.DB.Admin.Password = utils.GetRandomBase64(30)
	}
}

func (flags *UpgradeFlags) SetPasswordIfMissing() {
	flags.DB.SetPasswordIfMissing()
	flags.ReportDB.SetPasswordIfMissing()

	// The admin password is only needed for local database
	if flags.DB.IsLocal() && flags.DB.Admin.Password == "" {
		flags.DB.Admin.Password = utils.GetRandomBase64(30)
	}
}

func (flags *InstallationFlags) CheckSSLParameters(cmd *cobra.Command, command string) {
	// Make sure we have all the required 3rd party flags or none
	flags.SSL.CheckParameters(flags.DB.IsLocal())

	// Since we use cert-manager for self-signed certificates on kubernetes we don't need password for it
	if !flags.SSL.Server.UseProvided() && command == "podman" {
		utils.AskPasswordIfMissing(&flags.SSL.Password, cmd.Flag("ssl-password").Usage, 0, 0)
	}
}

func (flags *UpgradeFlags) CheckSSLParameters(cmd *cobra.Command, command string) {
	isLocalDB := flags.DB.Host == "db"
	// Make sure we have all the required 3rd party flags or none
	flags.SSL.DB.CheckParameters(isLocalDB)

	// Since we use cert-manager for self-signed certificates on kubernetes we don't need password for it
	if !flags.SSL.DB.UseProvided() && command == "podman" {
		utils.AskPasswordIfMissing(&flags.SSL.Password, cmd.Flag("ssl-password").Usage, 0, 0)
	}
}

// IsLocal indicates if the database is a local or a third party one.
func (flags *DBFlags) IsLocal() bool {
	return flags.Host == "" || flags.Host == "db" || flags.Host == "reportdb"
}

// DebugFlags contains information about enabled/disabled debug.
type DebugFlags struct {
	Java bool
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

// SSHFlags is the structure holding the SSH configuration to use to connect to the source server to migrate.
type SSHFlags struct {
	Key struct {
		Public  string
		Private string
	}
	Knownhosts string
	Config     string
}
