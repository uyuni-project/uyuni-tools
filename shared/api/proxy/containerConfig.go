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
func ContainerConfig(client *api.APIClient, request ProxyConfigRequest) (*[]int8, error) {
	return executeRequest(client, ProxyConfigRequestToMap(request))
}

// Compute and download the configuration file for proxy containers.
func ContainerConfigGenerate(client *api.APIClient, request ProxyConfigGenerateRequest) (*[]int8, error) {
	return executeRequest(client, ProxyConfigGenerateRequestToMap(request))
}

// common method to execute the request.
func executeRequest(client *api.APIClient, data map[string]interface{}) (*[]int8, error) {
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
