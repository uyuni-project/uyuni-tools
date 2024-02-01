// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

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

type SccFlags struct {
	User     string
	Password string
}

type DebugFlags struct {
	Java bool
}

type InstallFlags struct {
	TZ         string
	Email      string
	EmailFrom  string
	IssParent  string
	MirrorPath string
	Tftp       bool
	Db         DbFlags
	ReportDb   DbFlags
	Ssl        cmd_utils.SslCertFlags
	Scc        SccFlags
	Debug      DebugFlags
	Image      types.ImageFlags `mapstructure:",squash"`
}

func (flags *InstallFlags) CheckParameters(cmd *cobra.Command, command string) {
	utils.AskPasswordIfMissing(&flags.Db.Password, cmd.Flag("db-password").Usage)

	// Make sure we have all the required 3rd party flags or none
	flags.Ssl.CheckParameters()

	// Since we use cert-manager for self-signed certificates on kubernetes we don't need password for it
	if !flags.Ssl.UseExisting() && command == "podman" {
		utils.AskPasswordIfMissing(&flags.Ssl.Password, cmd.Flag("ssl-password").Usage)
	}

	// Use the host timezone if the user didn't define one
	if flags.TZ == "" {
		flags.TZ = utils.GetLocalTimezone()
	}

	utils.AskIfMissing(&flags.Email, cmd.Flag("email").Usage)
	utils.AskIfMissing(&flags.EmailFrom, cmd.Flag("emailfrom").Usage)
}

func AddInstallFlags(cmd *cobra.Command) {
	cmd.Flags().String("tz", "", "Time zone to set on the server. Defaults to the host timezone")
	cmd.Flags().String("email", "admin@example.com", "Administrator e-mail")
	cmd.Flags().String("emailfrom", "admin@example.com", "E-Mail sending the notifications")
	cmd.Flags().String("mirrorPath", "", "Path to mirrored packages mounted on the host")
	cmd.Flags().String("issParent", "", "Inter Server Sync v1 parent fully qualified domain name")
	cmd.Flags().String("db-user", "spacewalk", "Database user")
	cmd.Flags().String("db-password", "", "Database password")
	cmd.Flags().String("db-name", "susemanager", "Database name")
	cmd.Flags().String("db-host", "localhost", "Database host")
	cmd.Flags().Int("db-port", 5432, "Database port")
	cmd.Flags().String("db-protocol", "tcp", "Database protocol")
	cmd.Flags().String("db-admin-user", "", "External database admin user name")
	cmd.Flags().String("db-admin-password", "", "External database admin password")
	cmd.Flags().String("db-provider", "", "External database provider. Possible values 'aws'")

	cmd.Flags().Bool("tftp", true, "Enable TFTP")
	cmd.Flags().String("reportdb-name", "reportdb", "Report database name")
	cmd.Flags().String("reportdb-host", "", "Report database host. Defaults to the selected FQDN")
	cmd.Flags().Int("reportdb-port", 5432, "Report database port")
	cmd.Flags().String("reportdb-user", "pythia_susemanager", "Report Database username")
	cmd.Flags().String("reportdb-password", "", "Report database password. Randomly generated by default")

	// For generated CA and certificate
	cmd.Flags().StringSlice("ssl-cname", []string{}, "SSL certificate cnames separated by commas")
	cmd.Flags().String("ssl-country", "DE", "SSL certificate country")
	cmd.Flags().String("ssl-state", "Bayern", "SSL certificate state")
	cmd.Flags().String("ssl-city", "Nuernberg", "SSL certificate city")
	cmd.Flags().String("ssl-org", "SUSE", "SSL certificate organization")
	cmd.Flags().String("ssl-ou", "SUSE", "SSL certificate organization unit")
	cmd.Flags().String("ssl-password", "", "Password for the CA certificate to generate")
	cmd.Flags().String("ssl-email", "ca-admin@example.com", "SSL certificate E-Mail")

	// For SSL 3rd party certificates
	cmd.Flags().StringSlice("ssl-ca-intermediate", []string{}, "Intermediate CA certificate path")
	cmd.Flags().String("ssl-ca-root", "", "Root CA certificate path")
	cmd.Flags().String("ssl-server-cert", "", "Server certificate path")
	cmd.Flags().String("ssl-server-key", "", "Server key path")

	cmd.Flags().String("scc-user", "", "SUSE Customer Center username")
	cmd.Flags().String("scc-password", "", "SUSE Customer Center password")

	cmd.Flags().Bool("debug-java", false, "Enable tomcat and taskomatic remote debugging")
	cmd_utils.AddImageFlag(cmd)
}
