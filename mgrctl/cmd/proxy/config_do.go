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
	flags *proxyCreateConfigFlags,
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
	flags *proxyCreateConfigFlags,
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
