// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package hub

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupHubXmlrpc prepares the systemd service and starts it if needed.
// tag is the global images tag.
func SetupHubXmlrpc(
	authFile string,
	registry string,
	pullPolicy string,
	tag string,
	hubXmlrpcFlags cmd_utils.HubXmlrpcFlags,
) error {
	image := hubXmlrpcFlags.Image
	currentReplicas := podman.CurrentReplicaCount(podman.HubXmlrpcService)
	log.Debug().Msgf("Current HUB replicas running are %d.", currentReplicas)

	if hubXmlrpcFlags.Replicas == 0 {
		log.Debug().Msg("No HUB requested.")
	}
	if !hubXmlrpcFlags.IsChanged {
		log.Info().Msgf(L("No changes requested for hub. Keep %d replicas."), currentReplicas)
	}

	hubXmlrpcImage, err := utils.ComputeImage(registry, tag, image)

	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, hubXmlrpcImage, pullPolicy)
	if err != nil && ((hubXmlrpcFlags.Replicas > 0 && hubXmlrpcFlags.IsChanged) || (currentReplicas >= 0 && !hubXmlrpcFlags.IsChanged)) {
		return err
	} else if err != nil && ((hubXmlrpcFlags.Replicas == 0 && hubXmlrpcFlags.IsChanged) || (currentReplicas == 0 && !hubXmlrpcFlags.IsChanged)) {
		log.Info().Msgf(L("Image %s not present and it will not be pulled since HUB is not requested."), hubXmlrpcImage)
	}

	if err := generateHubXmlrpcSystemdService(preparedImage); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
	}

	if err := EnableHubXmlrpc(hubXmlrpcFlags.Replicas); err != nil {
		return err
	}
	return nil
}

// EnableHubXmlrpc enables the hub xmlrpc service if the number of replicas is 1.
// This function is meant for installation or migration, to enable or disable the service after, use ScaleService.
func EnableHubXmlrpc(replicas int) error {
	if replicas > 1 {
		log.Warn().Msg(L("Multiple Hub XML-RPC container replicas are not currently supported, setting up only one."))
		replicas = 1
	}

	if replicas > 0 {
		if err := podman.ScaleService(replicas, podman.HubXmlrpcService); err != nil {
			return utils.Errorf(err, L("cannot enable service"))
		}
	} else {
		log.Info().Msg(L("Not starting Hub XML-RPC API service"))
	}
	return nil
}

// Upgrade updates the systemd service files and restarts the containers if needed.
func Upgrade(authFile string, registry string, pullPolicy string, tag string, hubXmlrpcFlags cmd_utils.HubXmlrpcFlags) error {
	if err := SetupHubXmlrpc(authFile, registry, pullPolicy, tag, hubXmlrpcFlags); err != nil {
		return err
	}

	if err := podman.ReloadDaemon(false); err != nil {
		return err
	}

	return podman.RestartInstantiated(podman.HubXmlrpcService)
}

// generateHubXmlrpcSystemdService creates the Hub XMLRPC systemd files.
func generateHubXmlrpcSystemdService(image string) error {
	hubXmlrpcData := templates.HubXmlrpcServiceTemplateData{
		Volumes:    utils.HubXmlrpcVolumeMounts,
		Ports:      utils.HUB_XMLRPC_PORTS,
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Image:      image,
	}
	if err := utils.WriteTemplateToFile(hubXmlrpcData, podman.GetServicePath(podman.HubXmlrpcService+"@"), 0555, true); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf("Environment=UYUNI_IMAGE=%s", image)
	if err := podman.GenerateSystemdConfFile(podman.HubXmlrpcService+"@", "generated.conf", environment, true); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	return podman.ReloadDaemon(false)
}
