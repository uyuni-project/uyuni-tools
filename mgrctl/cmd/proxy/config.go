// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Specific flag names for proxy create config command.
const (
	proxyName      = "proxy-name"
	proxyPort      = "proxy-sshPort"
	server         = "proxy-parent"
	maxCache       = "proxy-maxCache"
	email          = "proxy-email"
	output         = "output"
	caCrt          = "ssl-ca-cert"
	caKey          = "ssl-ca-key"
	caPassword     = "ssl-ca-password"
	caIntermediate = "ssl-ca-intermediate"
	proxyCrt       = "ssl-proxy-cert"
	proxyKey       = "ssl-proxy-key"
	sslEmail       = "ssl-email"
)

type proxyFlags struct {
	Name     string
	Port     int `mapstructure:"sshPort"`
	Parent   string
	MaxCache int
	Email    string
}

type caFlags struct {
	types.SslPair `mapstructure:",squash"`
	Password      string
	Intermediate  []string
}

type proxyConfigSslFlags struct {
	types.SslCertGenerationFlags `mapstructure:",squash"`
	Proxy                        types.SslPair
	Ca                           caFlags
}

// proxyCreateConfigFlags is the structure containing the flags for proxy create config command.
type proxyCreateConfigFlags struct {
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
	Proxy             proxyFlags
	Output            string
	Ssl               proxyConfigSslFlags
}

// proxyCreateConfigRequiredFields is a set of required fields for validation.
var proxyCreateConfigRequiredFields = []string{
	proxyName,
	server,
	email,
	caCrt,
}

func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[proxyCreateConfigFlags]) *cobra.Command {
	createConfigCmd := &cobra.Command{
		Use:   "config",
		Short: L("Create a proxy configuration file"),
		Long:  L("Create a proxy configuration file"),
		Example: `  Create a proxy configuration file providing certificates providing only required parameters

    $ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@example.com" --ssl-ca-cert="root_ca.pem" --ssl-proxy-cert="proxy_crt.pem" \
		--ssl-proxy-key="proxy_key.pem"

  Create a proxy configuration file providing certificates providing  all parameters

    $ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@example.com" --ssl-ca-cert="root_ca.pem" --ssl-proxy-cert="proxy_crt.pem" \
		--ssl-proxy-key="proxy_key.pem" \
		--ssl-ca-intermediate="intermediateCA_1.pem,intermediateCA_2.pem,intermediateCA_3.pem" \
		-o="proxy-config"

  or an alternative format:

    $ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@example.com" --ssl-ca-cert="root_ca.pem" --ssl-proxy-cert="proxy_crt.pem" \
		--ssl-proxy-key="proxy_key.pem" \
		--ssl-ca-intermediate "intermediateCA_1.pem" --ssl-ca-intermediate "intermediateCA_2.pem" \
		--ssl-ca-intermediate "intermediateCA_3.pem" -o="proxy-config"

  Create a proxy configuration file with generated certificates providing only required parameters

    $ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@org.com" --ssl-ca-cert="ca.pem" --ssl-ca-key="caKey.pem"

  Create a proxy configuration file with generated certificates providing only required parameters and ca password

    $ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@org.com" --ssl-ca-cert="ca.pem" --ssl-ca-key="caKey.pem" --ssl-ca-password="secret"

  Create a proxy configuration file with generated certificates providing all parameters

    $ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@org.com" --ssl-ca-cert="ca.pem" --ssl-ca-key="caKey.pem" --ssl-ca-password="secret" \
		--ssl-cnames="proxy_a.example.com,proxy_b.example.com,proxy_c.example.com" --ssl-country="DE" \
		--ssl-state="Bayern" --ssl-city="Nuernberg" --ssl-org="orgExample" --ssl-ou="orgUnitExample" \
		--ssl-email="sslEmail@example.com" -o="proxy-config"

  or an alternative format:

	$ mgrctl proxy create config --proxy-name="proxy.example.com" --proxy-parent="server.example.com" \
		--proxy-email="admin@org.com" --ssl-ca-cert="ca.pem" --ssl-ca-key="caKey.pem" --ssl-ca-password="secret" \
		--ssl-cnames="proxy_a.example.com" --ssl-cnames="proxy_b.example.com" --ssl-cnames="proxy_c.example.com" \
		--ssl-country="DE" --ssl-state="Bayern" --ssl-city="Nuernberg" --ssl-org="orgExample" --ssl-ou="orgUnitExample" \
		--ssl-email="sslEmail@example.com" -o="proxy-config"

  Note that passing the CA password using --ssl-ca-password is not secure, use --config config.yaml with config.yaml
  containing the following as this will not persist in the shell history. Alternativaly the password can be defined in
  an UYUNI_SSL_CA_PASSWORD environment variable.

  ssl:
    ca:
	  password: secret
  `,
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags proxyCreateConfigFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
	}

	addFlags(createConfigCmd)

	// validations
	utils.MarkMandatoryFlags(createConfigCmd, proxyCreateConfigRequiredFields[:])
	createConfigCmd.MarkFlagsOneRequired(proxyCrt, caKey)
	createConfigCmd.MarkFlagsMutuallyExclusive(proxyCrt, caKey)

	return createConfigCmd
}

// NewConfigCommand creates the command for managing cache.
// Setup for subcommand to clear (the cache).
func NewConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, proxyCreateConfigInit)
}

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(output, "o", "", L("Filename to write the configuration to (without extension)."))

	// Common flags in command scope
	cmd.Flags().String(proxyName, "", L("Unique DNS-resolvable FQDN of this proxy."))
	cmd.Flags().Int(proxyPort, 8022, L("SSH port the proxy listens on."))
	cmd.Flags().String(server, "", L("FQDN of the server to connect the proxy to."))
	cmd.Flags().Int(maxCache, 102400, L("Maximum cache size in MB."))
	cmd.Flags().String(email, "", L("Email of the proxy administrator"))
	cmd.Flags().String(caCrt, "", L("Path to the root CA certificate in PEM format."))

	// Specific flags for when providing proxy certificates
	cmd.Flags().String(proxyCrt, "", L("Path to the proxy certificate in PEM format."))
	cmd.Flags().String(proxyKey, "", L("Path to the proxy certificate private key in PEM format."))
	cmd.Flags().StringSliceP(caIntermediate, "i", []string{},
		L(`Path to an intermediate CA used to sign the proxy certicate in PEM format.
May be provided multiple times or separated by commas.`),
	)

	// Specific flags for when generating proxy certificates
	ssl.AddSSLGenerationFlags(cmd)
	cmd.Flags().String(sslEmail, "", L("Email to set in the SSL certificate"))

	cmd.Flags().String(caKey, "", L("Path to the private key of the CA to use to generate a new proxy certificate."))
	cmd.Flags().String(caPassword, "",
		L("Password of the CA private key, will be prompted if not passed."),
	)

	// Login API flags
	api.AddAPIFlags(cmd)

	// Setup flag groups
	commonGroup := "common"
	providingGroup := "providing"
	_ = utils.AddFlagHelpGroup(cmd,
		&utils.Group{ID: commonGroup, Title: L("Common Flags")},
		&utils.Group{ID: providingGroup, Title: L("Third party proxy certificates flags")},
	)

	_ = utils.AddFlagsToHelpGroupID(cmd, commonGroup, proxyName, proxyPort, server, maxCache, email, caCrt)
	_ = utils.AddFlagsToHelpGroupID(cmd, providingGroup, proxyCrt, proxyKey, caIntermediate)

	_ = utils.AddFlagsToHelpGroupID(cmd, "ssl", sslEmail)
}
