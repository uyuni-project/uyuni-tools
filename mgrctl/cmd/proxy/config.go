// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/proxy"
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

// ProxyCreateConfigFlags is the structure containing the flags for proxy create config command.
type ProxyCreateConfigFlags struct {
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

// ProxyCreateConfigRequiredFields is a set of required fields for validation.
var ProxyCreateConfigRequiredFields = [6]string{
	proxyName,
	server,
	email,
	caCrt,
}

// NewConfigCommand creates the command for managing cache.
// Setup for subcommand to clear (the cache).
func NewConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags ProxyCreateConfigFlags

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
	utils.MarkMandatoryFlags(createConfigCmd, ProxyCreateConfigRequiredFields[:])
	createConfigCmd.MarkFlagsOneRequired(proxyCrt, caKey)
	createConfigCmd.MarkFlagsMutuallyExclusive(proxyCrt, caKey)

	return createConfigCmd
}

func proxyCreateConfigInit(
	globalFlags *types.GlobalFlags,
	flags *ProxyCreateConfigFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return ProxyCreateConfig(flags, api.Init, proxy.ContainerConfig, proxy.ContainerConfigGenerate)
}

// ProxyCreateConfig command handler.
func ProxyCreateConfig(
	flags *ProxyCreateConfigFlags,
	apiInit func(*api.ConnectionDetails) (*api.APIClient, error),
	proxyConfig func(client *api.APIClient, request proxy.ProxyConfigRequest) (*[]int8, error),
	proxyConfigGenerate func(client *api.APIClient, request proxy.ProxyConfigGenerateRequest) (*[]int8, error),
) error {
	client, err := apiInit(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}

	if err != nil {
		return utils.Errorf(err, L("failed to connect to the server"))
	}

	// handle CA certificate path
	caCertificate := string(utils.ReadFile(flags.CaCrt))

	// Check if ProxyCrt is provided to decide which configuration to run
	var data *[]int8
	if flags.ProxyCrt != "" {
		data, err = handleProxyConfig(client, flags, caCertificate, proxyConfig)
	} else {
		data, err = handleProxyConfigGenerate(client, flags, caCertificate, proxyConfigGenerate)
	}

	if err != nil {
		return utils.Errorf(err, L("failed to execute proxy configuration api request"))
	}

	filename := GetFilename(flags.Output, flags.ProxyName)
	if err := utils.SaveBinaryData(filename, *data); err != nil {
		return utils.Errorf(err, L("error saving binary data: %v"), err)
	}
	log.Info().Msgf(L("Proxy configuration file saved as %s"), filename)

	return nil
}

// Helper function to handle proxy configuration.
func handleProxyConfig(
	client *api.APIClient,
	flags *ProxyCreateConfigFlags,
	caCertificate string,
	proxyConfig func(client *api.APIClient, request proxy.ProxyConfigRequest) (*[]int8, error),
) (*[]int8, error) {
	// Custom validations
	if flags.ProxyKey == "" {
		return nil, errors.New(L("flag proxyKey is required when flag proxyCrt is provided"))
	}

	// Read file paths for certificates and keys
	proxyCrt := string(utils.ReadFile(flags.ProxyCrt))
	proxyKey := string(utils.ReadFile(flags.ProxyKey))

	// Handle intermediate CAs
	var intermediateCAs []string
	for _, path := range flags.IntermediateCAs {
		intermediateCAs = append(intermediateCAs, string(utils.ReadFile(path)))
	}

	// Prepare the request object & call the proxyConfig function
	request := proxy.ProxyConfigRequest{
		ProxyName:       flags.ProxyName,
		ProxyPort:       flags.ProxyPort,
		Server:          flags.Server,
		MaxCache:        flags.MaxCache,
		Email:           flags.Email,
		RootCA:          caCertificate,
		ProxyCrt:        proxyCrt,
		ProxyKey:        proxyKey,
		IntermediateCAs: intermediateCAs,
	}

	return proxyConfig(client, request)
}

// Helper function to handle proxy configuration generation.
func handleProxyConfigGenerate(
	client *api.APIClient,
	flags *ProxyCreateConfigFlags,
	caCertificate string,
	proxyConfigGenerate func(client *api.APIClient, request proxy.ProxyConfigGenerateRequest) (*[]int8, error),
) (*[]int8, error) {
	// CA key and password
	caKey := string(utils.ReadFile(flags.CaKey))

	var caPasswordRead string
	if flags.CaPassword == "" {
		utils.AskPasswordIfMissingOnce(&caPasswordRead, L("Please enter "+caPassword), 0, 0)
	} else {
		caPasswordRead = string(utils.ReadFile(flags.CaPassword))
	}

	// Prepare the request object & call the proxyConfigGenerate function
	request := proxy.ProxyConfigGenerateRequest{
		ProxyName:  flags.ProxyName,
		ProxyPort:  flags.ProxyPort,
		Server:     flags.Server,
		MaxCache:   flags.MaxCache,
		Email:      flags.Email,
		CaCrt:      caCertificate,
		CaKey:      caKey,
		CaPassword: caPasswordRead,
		Cnames:     flags.CNames,
		Country:    flags.Country,
		State:      flags.State,
		City:       flags.City,
		Org:        flags.Org,
		OrgUnit:    flags.OrgUnit,
		SslEmail:   flags.SslEmail,
	}

	return proxyConfigGenerate(client, request)
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
