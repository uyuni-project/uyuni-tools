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

func proxyCreateConfigInit(
	globalFlags *types.GlobalFlags,
	flags *proxyCreateConfigFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return proxyCreateConfig(flags, api.Init, proxy.ContainerConfig, proxy.ContainerConfigGenerate)
}

// proxyCreateConfig command handler.
func proxyCreateConfig(
	flags *proxyCreateConfigFlags,
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
	caCertificate := string(utils.ReadFile(flags.Ssl.Ca.Cert))

	// Check if ProxyCrt is provided to decide which configuration to run
	var data *[]int8
	if flags.Ssl.Proxy.Cert != "" {
		data, err = handleProxyConfig(client, flags, caCertificate, proxyConfig)
	} else {
		data, err = handleProxyConfigGenerate(client, flags, caCertificate, proxyConfigGenerate)
	}

	if err != nil {
		return utils.Errorf(err, L("failed to execute proxy configuration api request"))
	}

	filename := GetFilename(flags.Output, flags.Proxy.Name)
	if err := utils.SaveBinaryData(filename, *data); err != nil {
		return utils.Errorf(err, L("error saving binary data: %v"), err)
	}
	log.Info().Msgf(L("Proxy configuration file saved as %s"), filename)

	return nil
}

// Helper function to handle proxy configuration.
func handleProxyConfig(
	client *api.APIClient,
	flags *proxyCreateConfigFlags,
	caCertificate string,
	proxyConfig func(client *api.APIClient, request proxy.ProxyConfigRequest) (*[]int8, error),
) (*[]int8, error) {
	// Custom validations
	if flags.Ssl.Proxy.Key == "" {
		return nil, errors.New(L("flag proxyKey is required when flag proxyCrt is provided"))
	}

	// Read file paths for certificates and keys
	proxyCrt := string(utils.ReadFile(flags.Ssl.Proxy.Cert))
	proxyKey := string(utils.ReadFile(flags.Ssl.Proxy.Key))

	// Handle intermediate CAs
	var intermediateCAs []string
	for _, path := range flags.Ssl.Ca.Intermediate {
		intermediateCAs = append(intermediateCAs, string(utils.ReadFile(path)))
	}

	// Prepare the request object & call the proxyConfig function
	request := proxy.ProxyConfigRequest{
		ProxyName:       flags.Proxy.Name,
		ProxyPort:       flags.Proxy.Port,
		Server:          flags.Proxy.Parent,
		MaxCache:        flags.Proxy.MaxCache,
		Email:           flags.Proxy.Email,
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
	flags *proxyCreateConfigFlags,
	caCertificate string,
	proxyConfigGenerate func(client *api.APIClient, request proxy.ProxyConfigGenerateRequest) (*[]int8, error),
) (*[]int8, error) {
	// CA key and password
	caKey := string(utils.ReadFile(flags.Ssl.Ca.Key))

	var caPasswordRead string
	if flags.Ssl.Ca.Password == "" {
		utils.AskPasswordIfMissingOnce(&caPasswordRead, L("Please enter "+caPassword), 0, 0)
	} else {
		caPasswordRead = string(utils.ReadFile(flags.Ssl.Ca.Password))
	}

	// Prepare the request object & call the proxyConfigGenerate function
	request := proxy.ProxyConfigGenerateRequest{
		ProxyName:  flags.Proxy.Name,
		ProxyPort:  flags.Proxy.Port,
		Server:     flags.Proxy.Parent,
		MaxCache:   flags.Proxy.MaxCache,
		Email:      flags.Proxy.Email,
		CaCrt:      caCertificate,
		CaKey:      caKey,
		CaPassword: caPasswordRead,
		Cnames:     flags.Ssl.Cnames,
		Country:    flags.Ssl.Country,
		State:      flags.Ssl.State,
		City:       flags.Ssl.City,
		Org:        flags.Ssl.Org,
		OrgUnit:    flags.Ssl.OU,
		SslEmail:   flags.Ssl.Email,
	}

	return proxyConfigGenerate(client, request)
}
