// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
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
	rootCA          = "rootCA"
	intermediateCAs = "intermediateCAs"
	proxyCrt        = "proxyCrt"
	proxyKey        = "proxyKey"
)

// Flags for proxy create config command.
type ProxyCreateConfigFlags struct {
	ProxyCreateConfigBaseFlags `mapstructure:",squash"`
	RootCA                     string
	IntermediateCAs            []string
	ProxyCrt                   string
	ProxyKey                   string
}

// Set of required fields for validation.
var ProxyCreateConfigRequiredFields = [6]string{
	proxyName,
	server,
	email,
	rootCA,
	proxyCrt,
	proxyKey,
}

// CreateCommand entry command for managing cache.
// Setup for subcommand to clear (the cache).
func NewConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	var flags ProxyCreateConfigFlags

	createConfigCmd := &cobra.Command{
		Use:   "config",
		Short: L("Create a proxy configuration file"),
		Long:  L("Create a proxy configuration file"),
		Example: `  Create a proxy configuration file providing only required parameters
	
    $ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" --email="admin@example.com" --rootCA="root_ca.pem" --proxyCrt="proxy_crt.pem" --proxyKey="proxy_key.pem"

  Create a proxy configuration file providing all parameters
	
	$ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" --email="admin@example.com" --rootCA="root_ca.pem" --proxyCrt="proxy_crt.pem" --proxyKey="proxy_key.pem" --intermediateCAs="intermediateCA_1.pem,intermediateCA_2.pem,intermediateCA_3.pem" -o="proxy-config"
	
  or an alternative format:

	$ mgrctl proxy create config --proxyName="proxy.example.com" --server="server.example.com" --email="admin@example.com" --rootCA="root_ca.pem" --proxyCrt="proxy_crt.pem" --proxyKey="proxy_key.pem" --intermediateCAs "intermediateCA_1.pem" --intermediateCAs "intermediateCA_2.pem" --intermediateCAs "intermediateCA_3.pem" -o="proxy-config"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			viper, _ := utils.ReadConfig(cmd, utils.GlobalConfigFilename, globalFlags.ConfigPath)
			if err := viper.Unmarshal(&flags); err != nil {
				log.Fatal().Err(err).Msg(L("failed to unmarshall configuration"))
			}
			return ProxyCreateConfig(&flags, api.Init, proxy.ContainerConfig)
		},
	}

	createConfigCmd.Flags().String(proxyName, "", L("Unique DNS-resolvable FQDN of this proxy. Mandatory."))
	createConfigCmd.Flags().Int(proxyPort, 8022, L("SSH port the proxy listens one."))
	createConfigCmd.Flags().String(server, "", L("FQDN of the server to connect to proxy to. Mandatory."))
	createConfigCmd.Flags().Int(maxCache, 102400, L("Maximum cache size in MB."))
	createConfigCmd.Flags().String(email, "", L("Email of the proxy administrator"))
	createConfigCmd.Flags().StringP(output, "o", "", L("Filename to write the configuration to (without extension)."))

	createConfigCmd.Flags().String(rootCA, "", L("Path to the root CA used to sign the proxy certificate in PEM format. Mandatory."))
	createConfigCmd.Flags().String(proxyCrt, "", L("Path to the proxy certificate in PEM format. Mandatory."))
	createConfigCmd.Flags().String(proxyKey, "", L("Path to the proxy certificate private key in PEM format. Mandatory."))
	createConfigCmd.Flags().StringSliceP(intermediateCAs, "i", []string{}, L("Path to an intermediate CA used to sign the proxy certicate in PEM format. May be provided multiple times or separated by commas."))

	createConfigCmd.AddCommand(NewConfigGenerateCommand(globalFlags))

	utils.ValidateMandatoryFlags(createConfigCmd, ProxyCreateConfigRequiredFields[:])

	api.AddAPIFlags(createConfigCmd)
	return createConfigCmd
}

// ProxyCreateConfig command handler.
func ProxyCreateConfig(
	flags *ProxyCreateConfigFlags,
	apiInit func(*api.ConnectionDetails) (*api.APIClient, error),
	proxyConfig func(client *api.APIClient, proxyName string, proxyPort int,
		server string, maxCache int, email string,
		rootCA string, proxyCrt string, proxyKey string, intermediateCAs []string) (*[]int8, error),
) error {
	client, err := apiInit(&flags.ConnectionDetails)
	if err == nil {
		err = client.Login()
	}

	if err != nil {
		return utils.Errorf(err, L("failed to connect to the server"))
	}

	// handle file paths
	rootCA := string(utils.ReadFile(flags.RootCA))
	proxyCrt := string(utils.ReadFile(flags.ProxyCrt))
	proxyKey := string(utils.ReadFile(flags.ProxyKey))

	var intermediateCAs []string
	for _, v := range flags.IntermediateCAs {
		intermediateCAs = append(intermediateCAs, string(utils.ReadFile(v)))
	}

	data, err := proxyConfig(client, flags.ProxyName, flags.ProxyPort,
		flags.Server, flags.MaxCache, flags.Email,
		rootCA, proxyCrt, proxyKey, intermediateCAs)

	if err != nil {
		return utils.Errorf(err, L("failed to execute proxy configuration api request"))
	}

	filename := GetFilename(flags.Output, proxyName)
	if err := utils.SaveBinaryData(filename, *data); err != nil {
		return utils.Errorf(err, L("Error saving binary data: %v"), err)
	}

	log.Debug().Msgf("Proxy configuration file saved as %s", filename)

	return nil
}
