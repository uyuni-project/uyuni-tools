// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Specific flag names for proxy create config command.
const (
	proxyName       = "proxyName"
	proxyPort       = "proxyPort"
	server          = "server"
	maxCache        = "maxCache"
	email           = "email"
	output          = "output"
	caCrt           = "caCrt"
	caKey           = "caKey"
	caPassword      = "caPassword"
	cNames          = "cnames"
	country         = "country"
	state           = "state"
	city            = "city"
	org             = "org"
	orgUnit         = "orgUnit"
	sslEmail        = "sslEmail"
	intermediateCAs = "intermediateCAs"
	proxyCrt        = "proxyCrt"
	proxyKey        = "proxyKey"
)

// proxyCreateConfigFlags is the structure containing the flags for proxy create config command.
type proxyCreateConfigFlags struct {
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
	ProxyName         string
	ProxyPort         int
	Server            string
	MaxCache          int
	Email             string
	Output            string
	CaCrt             string
	ProxyCrt          string
	ProxyKey          string
	IntermediateCAs   []string
	CaKey             string
	CaPassword        string
	CNames            []string
	Country           string
	State             string
	City              string
	Org               string
	OrgUnit           string
	SslEmail          string
}

// proxyCreateConfigRequiredFields is a set of required fields for validation.
var proxyCreateConfigRequiredFields = [6]string{
	proxyName,
	server,
	email,
	caCrt,
}

// NewConfigCommand creates the command for managing cache.
// Setup for subcommand to clear (the cache).
func NewConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags proxyCreateConfigFlags

	createConfigCmd := &cobra.Command{
		Use:   "config",
		Short: L("Create a proxy configuration file"),
		Long:  L("Create a proxy configuration file"),
		Example: `  Create a proxy configuration file providing certificates providing only required parameters

    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@example.com" --caCrt="root_ca.pem" --proxyCrt="proxy_crt.pem" \
		--proxyKey="proxy_key.pem"

  Create a proxy configuration file providing certificates providing  all parameters

    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@example.com" --caCrt="root_ca.pem" --proxyCrt="proxy_crt.pem" \
		--proxyKey="proxy_key.pem" --intermediateCAs="intermediateCA_1.pem,intermediateCA_2.pem,intermediateCA_3.pem" \
		-o="proxy-config"

  or an alternative format:

    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@example.com" --caCrt="root_ca.pem" --proxyCrt="proxy_crt.pem" \
		--proxyKey="proxy_key.pem" --intermediateCAs "intermediateCA_1.pem" --intermediateCAs "intermediateCA_2.pem" \
		--intermediateCAs "intermediateCA_3.pem" -o="proxy-config"

  Create a proxy configuration file with generated certificates providing only required parameters

    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@org.com" --caCrt="ca.pem" --caKey="caKey.pem"

  Create a proxy configuration file with generated certificates providing only required parameters and ca password

    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@org.com" --caCrt="ca.pem" --caKey="caKey.pem" --caPassword="pass.txt"

  Create a proxy configuration file with generated certificates providing all parameters

    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@org.com" --caCrt="ca.pem" --caKey="caKey.pem" --caPassword="pass.txt" \
		--cnames="proxy_a.example.com,proxy_b.example.com,proxy_c.example.com" --country="DE" \
		--state="Bayern" --city="Nuernberg" --org="orgExample" --orgUnit="orgUnitExample" \
		--sslEmail="sslEmail@example.com" -o="proxy-config"

  or an alternative format:

	$ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" \
		--email="admin@org.com" --caCrt="ca.pem" --caKey="caKey.pem" --caPassword="pass.txt" \
		--cnames="proxy_a.example.com" --cnames="proxy_b.example.com" --cnames="proxy_c.example.com" \
		--country="DE" --state="Bayern" --city="Nuernberg" --org="orgExample" --orgUnit="orgUnitExample" \
		--sslEmail="sslEmail@example.com" -o="proxy-config"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, proxyCreateConfigInit)
		},
	}

	addFlags(createConfigCmd)

	// validations
	utils.MarkMandatoryFlags(createConfigCmd, proxyCreateConfigRequiredFields[:])
	createConfigCmd.MarkFlagsOneRequired(proxyCrt, caKey)
	createConfigCmd.MarkFlagsMutuallyExclusive(proxyCrt, caKey)

	return createConfigCmd
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
	cmd.Flags().StringSliceP(intermediateCAs, "i", []string{},
		L(`Path to an intermediate CA used to sign the proxy certicate in PEM format.
May be provided multiple times or separated by commas.`),
	)

	// Specific flags for when generating proxy certificates
	cmd.Flags().String(caKey, "", L("Path to the private key of the CA to use to generate a new proxy certificate."))
	cmd.Flags().String(caPassword, "",
		L("Path to a file containing the password of the CA private key, will be prompted if not passed."),
	)
	cmd.Flags().StringSlice(cNames, []string{},
		L("Alternate name of the proxy to set in the certificate. May be provided multiple times or separated by commas."),
	)
	cmd.Flags().String(country, "", L("Country code to set in the certificate."))
	cmd.Flags().String(state, "", L("State name to set in the certificate."))
	cmd.Flags().String(city, "", L("City name to set in the certificate."))
	cmd.Flags().String(org, "", L("Organization name to set in the certificate."))
	cmd.Flags().String(orgUnit, "", L("Organization unit name to set in the certificate."))
	cmd.Flags().String(sslEmail, "", L("Email to set in the certificate."))

	// Login API flags
	api.AddAPIFlags(cmd)

	// Setup flag groups
	commonGroup := "common"
	providingGroup := "providing"
	generateGroup := "generate"
	_ = utils.AddFlagHelpGroup(cmd,
		&utils.Group{ID: commonGroup, Title: L("Common Flags")},
		&utils.Group{ID: providingGroup, Title: L("Provide proxy certificates flags")},
		&utils.Group{ID: generateGroup, Title: L("Generate proxy certificates flags")},
	)

	_ = utils.AddFlagsToHelpGroupID(cmd, commonGroup, proxyName, proxyPort, server, maxCache, email, caCrt)
	_ = utils.AddFlagsToHelpGroupID(cmd, providingGroup, proxyCrt, proxyKey, intermediateCAs)
	_ = utils.AddFlagsToHelpGroupID(
		cmd, generateGroup, caKey, caPassword, cNames, country, state, city, org, orgUnit, sslEmail,
	)
}
