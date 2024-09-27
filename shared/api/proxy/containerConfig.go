// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const containerConfigEndpoint = "proxy/containerConfig"

// Compute and download the configuration file for proxy containers with generated certificates.
func ContainerConfig(client *api.APIClient, proxyName string, proxyPort int,
	server string, maxCache int, email string,
	rootCA string, proxyCrt string, proxyKey string, intermediateCAs []string) (*[]int8, error) {
	data := map[string]interface{}{
		"proxyName":       proxyName,
		"proxyPort":       proxyPort,
		"server":          server,
		"maxCache":        maxCache,
		"email":           email,
		"rootCA":          rootCA,
		"proxyCrt":        proxyCrt,
		"proxyKey":        proxyKey,
		"intermediateCAs": intermediateCAs,
	}
	log.Trace().Msgf("Creating proxy configuration file with generated certificates with data: %v...", data)

	res, err := api.Post[[]int8](client, containerConfigEndpoint, data)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to create proxy configuration file with generated certificates"))
	}
	if !res.Success {
		return nil, errors.New(res.Message)
	}

	return &res.Result, nil
}

// Compute and download the configuration file for proxy containers.
func ContainerConfigGenerate(client *api.APIClient, proxyName string, proxyPort int,
	server string, maxCache int, email string,
	caCertificate string, caKey string, caPassword string, cnames []string, country string,
	state string, city string, org string, orgUnit string, sslEmail string) (*[]int8, error) {
	data := map[string]interface{}{
		"proxyName":  proxyName,
		"proxyPort":  proxyPort,
		"server":     server,
		"maxCache":   maxCache,
		"email":      email,
		"caCrt":      caCertificate,
		"caKey":      caKey,
		"caPassword": caPassword,
		"cnames":     cnames,
		"country":    country,
		"state":      state,
		"city":       city,
		"org":        org,
		"orgUnit":    orgUnit,
		"sslEmail":   sslEmail,
	}
	log.Trace().Msgf("Creating proxy configuration file with data: %v...", data)

	res, err := api.Post[[]int8](client, containerConfigEndpoint, data)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to create proxy configuration file"))
	}
	if !res.Success {
		return nil, errors.New(res.Message)
	}
	return &res.Result, nil
}
