// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	"github.com/uyuni-project/uyuni-tools/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupHubXmlrpc prepares the systemd service and starts it if needed.
// tag is the global images tag.
func SetupHubXmlrpc(
	systemd podman.Systemd,
	authFile string,
	pullPolicy string,
	hubXmlrpcFlags cmd_utils.HubXmlrpcFlags,
) error {
	image := hubXmlrpcFlags.Image
	currentReplicas := systemd.CurrentReplicaCount(podman.HubXmlrpcService)
	log.Debug().Msgf("Current HUB replicas running are %d.", currentReplicas)

	if hubXmlrpcFlags.Replicas == 0 {
		log.Debug().Msg("No HUB requested.")
	}
	if !hubXmlrpcFlags.IsChanged && hubXmlrpcFlags.Replicas == currentReplicas {
		log.Info().Msgf(L("No changes requested for hub. Keep %d replicas."), currentReplicas)
	}

	pullEnabled := hubXmlrpcFlags.Replicas > 0 || (currentReplicas > 0 && !hubXmlrpcFlags.IsChanged)

	hubXmlrpcImage, err := utils.ComputeImage(image)

	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	preparedImage, err := podman.PrepareImage(authFile, hubXmlrpcImage, pullPolicy, pullEnabled)
	if err != nil {
		return err
	}

	if err := generateHubXmlrpcSystemdService(systemd, preparedImage.Name, podman.ServerContainerName); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
	}

	if err := EnableHubXmlrpc(systemd, hubXmlrpcFlags.Replicas); err != nil {
		return err
	}
	return nil
}

// EnableHubXmlrpc enables the hub xmlrpc service if the number of replicas is 1.
// This function is meant for installation or migration, to enable or disable the service after, use ScaleService.
func EnableHubXmlrpc(systemd podman.Systemd, replicas int) error {
	if replicas > 1 {
		log.Warn().Msg(L("Multiple Hub XML-RPC container replicas are not currently supported, setting up only one."))
		replicas = 1
	}

	if replicas > 0 {
		if err := systemd.ScaleService(replicas, podman.HubXmlrpcService); err != nil {
			return utils.Errorf(err, L("cannot enable service"))
		}
	}
	return nil
}

// Upgrade updates the systemd service files and restarts the containers if needed.
func Upgrade(
	systemd podman.Systemd,
	authFile string,
	pullPolicy string,
	hubXmlrpcFlags cmd_utils.HubXmlrpcFlags,
) error {
	if hubXmlrpcFlags.Image.Name == "" {
		// Don't touch the hub service in ptf if not already present.
		return nil
	}
	if err := SetupHubXmlrpc(systemd, authFile, pullPolicy, hubXmlrpcFlags); err != nil {
		return err
	}

	if err := systemd.ReloadDaemon(false); err != nil {
		return err
	}

	if !hubXmlrpcFlags.IsChanged {
		return systemd.RestartInstantiated(podman.HubXmlrpcService)
	}
	return systemd.ScaleService(hubXmlrpcFlags.Replicas, podman.HubXmlrpcService)
}

// generateHubXmlrpcSystemdService creates the Hub XMLRPC systemd files.
func generateHubXmlrpcSystemdService(systemd podman.Systemd, image string, serverHost string) error {
	hubXmlrpcData := templates.HubXmlrpcServiceTemplateData{
		CaSecret:   podman.CASecret,
		CaPath:     ssl.CAContainerPath,
		Ports:      utils.HubXmlrpcPorts,
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Image:      image,
		ServerHost: serverHost,
	}
	if err := utils.WriteTemplateToFile(
		hubXmlrpcData, podman.GetServicePath(podman.HubXmlrpcService+"@"), 0555, true,
	); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf("Environment=UYUNI_HUB_XMLRPC_IMAGE=%s", image)
	if err := podman.GenerateSystemdConfFile(
		podman.HubXmlrpcService+"@", "generated.conf", environment, true,
	); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	return systemd.ReloadDaemon(false)
}
