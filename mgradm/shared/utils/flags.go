// SPDX-FileCopyrightText: 2024 SUSE LLC
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
	Image        types.ImageFlags `mapstructure:",squash"`
	Coco         CocoFlags
	Mirror       string
	HubXmlrpc    HubXmlrpcFlags
	Migration    MigrationFlags    `mapstructure:",squash"`
	Installation InstallationFlags `mapstructure:",squash"`
	// DBUpgradeImage is the image to use to perform the database upgrade.
	DBUpgradeImage types.ImageFlags `mapstructure:"dbupgrade"`
	Saline         SalineFlags
}

// MigrationFlags contains the parameters that are used only for migration.
type MigrationFlags struct {
	// Prepare defines whether to run the full migration or just the data synchronization.
	Prepare bool
	// SourceUser is the username to use to connect to the source server in a migration.
	User string
}

// InstallationFlags contains the parameters that are used only for the installation of a new server.
type InstallationFlags struct {
	TZ           string
	Email        string
	EmailFrom    string
	IssParent    string
	Tftp         bool
	DB           DBFlags
	ReportDB     DBFlags
	SSL          InstallSSLFlags
	SCC          types.SCCCredentials
	Debug        DebugFlags
	Admin        apiTypes.User
	Organization string
}

// CheckParameters checks parameters for install command.
func (flags *InstallationFlags) CheckParameters(cmd *cobra.Command, command string) {
	if flags.DB.Password == "" {
		flags.DB.Password = utils.GetRandomBase64(30)
	}

	if flags.ReportDB.Password == "" {
		flags.ReportDB.Password = utils.GetRandomBase64(30)
	}

	// Make sure we have all the required 3rd party flags or none
	flags.SSL.CheckParameters()

	// Since we use cert-manager for self-signed certificates on kubernetes we don't need password for it
	if !flags.SSL.UseExisting() && command == "podman" {
		utils.AskPasswordIfMissing(&flags.SSL.Password, cmd.Flag("ssl-password").Usage, 0, 0)
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

	flags.SSL.Email = flags.Email
	flags.Admin.Email = flags.Email
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
