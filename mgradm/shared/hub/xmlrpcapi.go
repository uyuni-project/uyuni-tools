// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package hub

import (
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupHubXmlrpcContainer prepares the systemd service and starts it if needed.
// tag is the global images tag.
func HubXmlrpc(tag string, hubxmlrpcFlags *shared.HubXmlrpcFlags) error {
	if hubxmlrpcFlags.Replicas > 1 {
		log.Warn().Msg(L("Multiple Hub XML-RPC container replicas are not currently supported, setting up only one."))
		hubxmlrpcFlags.Replicas = 1
	}
	log.Info().Msg(L("Setting Hub XML-RPC API service."))
	if hubxmlrpcFlags.Image.Tag == "" {
		hubxmlrpcFlags.Image.Tag = tag
	}
	hubXmlrpcImage, err := utils.ComputeImage(hubxmlrpcFlags.Image)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	if err := podman.GenerateHubXmlrpcSystemdService(hubXmlrpcImage); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
	}

	if hubxmlrpcFlags.Replicas > 0 {
		if err := shared_podman.ScaleService(hubxmlrpcFlags.Replicas, shared_podman.HubXmlrpcService); err != nil {
			return utils.Errorf(err, L("cannot enable service"))
		}
	} else {
		log.Info().Msg(L("Not starting Hub XML-RPC API service"))
	}
	return nil
}
