// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/mgradm/shared/coco"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/hub"
	adm_podman "github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/saline"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

var systemd podman.Systemd = podman.NewSystemd()

func podmanScale(
	_ *types.GlobalFlags,
	flags *scaleFlags,
	_ *cobra.Command,
	args []string,
) error {
	newReplicas := flags.Replicas
	service := args[0]
	if service == podman.ServerAttestationService {
		return scaleService(flags, service)
	}
	if service == podman.HubXmlrpcService || service == podman.SalineService {
		if newReplicas > 1 {
			return errors.New(L("Multiple container replicas are not currently supported."))
		}
		return scaleService(flags, service)
	}
	return fmt.Errorf(L("service not allowing to be scaled: %s"), service)
}

// scaleService runs the upgrade code of the service, which pulls the image with the registry credentials.
func scaleService(flags *scaleFlags, service string) error {
	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman.PodmanLogin(hostData, flags.Image.Registry, flags.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	// Keep the current image of the service.
	image := types.ImageFlags{Name: podman.GetServiceImage(service + "@")}

	switch service {
	case podman.ServerAttestationService:
		cocoFlags := adm_utils.CocoFlags{Replicas: flags.Replicas, Image: image, IsChanged: true}
		return coco.Upgrade(systemd, authFile, cocoFlags, flags.Image)
	case podman.HubXmlrpcService:
		hubFlags := adm_utils.HubXmlrpcFlags{Replicas: flags.Replicas, Image: image, IsChanged: true}
		return hub.Upgrade(systemd, authFile, flags.Image, hubFlags)
	default:
		salineFlags := adm_utils.SalineFlags{Replicas: flags.Replicas, Image: image, IsChanged: true}
		return saline.Upgrade(systemd, authFile, flags.Image, salineFlags, adm_podman.GetContainerTimezone())
	}
}
