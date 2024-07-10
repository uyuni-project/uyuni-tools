// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package hub

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/templates"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// SetupHubXmlrpc prepares the systemd service and starts it if needed.
// tag is the global images tag.
func SetupHubXmlrpc(registry string, pullPolicy string, tag string, hubxmlrpcImage types.ImageFlags) error {
	log.Info().Msg(L("Setting Hub XML-RPC API service."))
	hubXmlrpcImage, err := utils.ComputeImage(registry, tag, hubxmlrpcImage)
	if err != nil {
		return utils.Errorf(err, L("failed to compute image URL"))
	}

	inspectedHostValues, err := utils.InspectHost(false)
	if err != nil {
		return utils.Errorf(err, L("cannot inspect host values"))
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := podman.PrepareImage(hubXmlrpcImage, pullPolicy, pullArgs...)
	if err != nil {
		return err
	}

	if err := generateHubXmlrpcSystemdService(preparedImage); err != nil {
		return utils.Errorf(err, L("cannot generate systemd service"))
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

// generateHubXmlrpcSystemdService creates the Hub XMLRPC systemd files.
func generateHubXmlrpcSystemdService(image string) error {
	hubXmlrpcData := templates.HubXmlrpcServiceTemplateData{
		Volumes:    utils.HubXmlrpcVolumeMounts,
		Ports:      utils.HUB_XMLRPC_PORTS,
		NamePrefix: "uyuni",
		Network:    podman.UyuniNetwork,
		Image:      image,
	}
	if err := utils.WriteTemplateToFile(hubXmlrpcData, podman.GetServicePath(podman.HubXmlrpcService+"@"), 0555, false); err != nil {
		return utils.Errorf(err, L("failed to generate systemd service unit file"))
	}

	environment := fmt.Sprintf(`Environment=UYUNI_IMAGE=%s
	`, image)
	if err := podman.GenerateSystemdConfFile(podman.HubXmlrpcService+"@", "Service", environment); err != nil {
		return utils.Errorf(err, L("cannot generate systemd conf file"))
	}

	return podman.ReloadDaemon(false)
}
