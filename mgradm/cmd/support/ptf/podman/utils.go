// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	podman_shared "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman_shared.Systemd = podman_shared.SystemdImpl{}

func ptfForPodman(
	globalFlags *types.GlobalFlags,
	flags *podmanPTFFlags,
	cmd *cobra.Command,
	args []string,
) error {
	// Login first to be able to search the registry for PTF images
	hostData, err := podman_shared.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman_shared.PodmanLogin(hostData, flags.SCC)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	//we don't want to perform a postgres version upgrade when installing a PTF.
	//in that case, we can use the upgrade command.
	dummyImage := types.ImageFlags{}
	if err := flags.checkParameters(); err != nil {
		return err
	}

	return podman.Upgrade(systemd, authFile, "", flags.Image, dummyImage, flags.Coco, flags.Hubxmlrpc)
}

// variables for unit testing.
var getServiceImage = podman_shared.GetServiceImage
var hasRemoteImage = podman_shared.HasRemoteImage

func (flags *podmanPTFFlags) checkParameters() error {
	if flags.TestID != "" && flags.PTFId != "" {
		return errors.New(L("ptf and test flags cannot be set simultaneously "))
	}
	if flags.TestID == "" && flags.PTFId == "" {
		return errors.New(L("ptf and test flags cannot be empty simultaneously "))
	}
	if flags.CustomerID == "" {
		return errors.New(L("user flag cannot be empty"))
	}

	suffix := "ptf"
	projectID := flags.PTFId
	if flags.TestID != "" {
		suffix = "test"
		projectID = flags.TestID
	}

	serverImage := getServiceImage(podman_shared.ServerService)
	if serverImage == "" {
		return errors.New(L("failed to find server image"))
	}

	var err error

	flags.Image.Name, err = utils.ComputePTF(flags.CustomerID, projectID, serverImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The computed image is %s"), flags.Image.Name)

	if cocoImage := getServiceImage(podman_shared.ServerAttestationService + "@"); cocoImage != "" {
		// If no coco image was found then skip it during the upgrade.
		cocoImage, err = utils.ComputePTF(flags.CustomerID, projectID, cocoImage, suffix)
		if err != nil {
			return err
		}
		if hasRemoteImage(cocoImage) {
			flags.Coco.Image.Name = cocoImage
			log.Info().Msgf(L("The computed confidential computing image is %s"), flags.Coco.Image.Name)
		}
	}

	if hubImage := getServiceImage(podman_shared.HubXmlrpcService); hubImage != "" {
		// If no hub xmlrpc api image was found then skip it during the upgrade.
		hubImage, err = utils.ComputePTF(flags.CustomerID, projectID, hubImage, suffix)
		if err != nil {
			return err
		}
		if hasRemoteImage(hubImage) {
			flags.Hubxmlrpc.Image.Name = hubImage
			log.Info().Msgf(L("The computed hub XML-RPC API image is %s"), flags.Hubxmlrpc.Image.Name)
		}
	}

	return nil
}
