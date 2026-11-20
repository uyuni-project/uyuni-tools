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
	"github.com/uyuni-project/uyuni-tools/shared/utils"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

var systemd podman.Systemd = podman.NewSystemd()

// variables for unit testing.
var getServiceImage = podman.GetServiceImage

func podmanScale(
	_ *types.GlobalFlags,
	flags *scaleFlags,
	cmd *cobra.Command,
	args []string,
) error {
	newReplicas := flags.Replicas
	service := args[0]
	if service == podman.ServerAttestationService {
		return scaleService(cmd, flags, service)
	}
	if service == podman.HubXmlrpcService || service == podman.SalineService {
		if newReplicas > 1 {
			return errors.New(L("Multiple container replicas are not currently supported."))
		}
		return scaleService(cmd, flags, service)
	}
	return fmt.Errorf(L("service not allowing to be scaled: %s"), service)
}

// scaleService runs the upgrade code of the service, which pulls the image with the registry credentials.
func scaleService(cmd *cobra.Command, flags *scaleFlags, service string) error {
	hostData, err := podman.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman.PodmanLogin(hostData, flags.Image.Registry, flags.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	// Same default tag as install and upgrade for an image passed without tag.
	if flags.Image.Tag == "" {
		flags.Image.Tag = utils.DefaultTag
	}

	switch service {
	case podman.ServerAttestationService:
		image := serviceImage(cmd, "coco", flags.Coco.Image, service)
		cocoFlags := adm_utils.CocoFlags{Replicas: flags.Replicas, Image: image, IsChanged: true}
		return coco.Upgrade(systemd, authFile, cocoFlags, flags.Image)
	case podman.HubXmlrpcService:
		image := serviceImage(cmd, "hubxmlrpc", flags.HubXmlrpc.Image, service)
		hubFlags := adm_utils.HubXmlrpcFlags{Replicas: flags.Replicas, Image: image, IsChanged: true}
		return hub.Upgrade(systemd, authFile, flags.Image, hubFlags)
	default:
		image := serviceImage(cmd, "saline", flags.Saline.Image, service)
		salineFlags := adm_utils.SalineFlags{Replicas: flags.Replicas, Image: image, IsChanged: true}
		return saline.Upgrade(systemd, authFile, flags.Image, salineFlags, adm_podman.GetContainerTimezone())
	}
}

// serviceImage returns the image passed with the <prefix>-image or <prefix>-tag flags, or the current one.
func serviceImage(cmd *cobra.Command, prefix string, image types.ImageFlags, service string) types.ImageFlags {
	if cmd.Flags().Changed(prefix+"-image") || cmd.Flags().Changed(prefix+"-tag") {
		return image
	}
	return types.ImageFlags{Name: getServiceImage(service + "@")}
}
