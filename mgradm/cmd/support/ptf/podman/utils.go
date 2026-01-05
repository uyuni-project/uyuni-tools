// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	podman_shared "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman_shared.Systemd = podman_shared.NewSystemd()

func ptfForPodman(
	_ *types.GlobalFlags,
	flags *podmanPTFFlags,
	_ *cobra.Command,
	_ []string,
) error {
	// Login first to be able to search the registry for PTF images
	hostData, err := podman_shared.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman_shared.PodmanLogin(hostData, flags.Image.Registry, flags.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	//we don't want to perform a postgres version upgrade when installing a PTF.
	//in that case, we can use the upgrade command.
	dummyImage := types.ImageFlags{}
	dummyDB := adm_utils.DBFlags{}
	dummyReportDB := adm_utils.DBFlags{}
	dummySSL := adm_utils.InstallSSLFlags{}

	if err := flags.checkParameters(authFile); err != nil {
		return err
	}

	return podman.Upgrade(systemd, authFile,
		dummyDB,
		dummyReportDB,
		dummySSL,
		flags.Image,
		dummyImage,
		flags.Coco,
		flags.HubXmlrpc,
		flags.Saline,
		flags.Pgsql,
		flags.Installation.TZ,
	)
}

// variables for unit testing.
var getServiceImage = podman_shared.GetServiceImage
var hasRemoteImage = podman_shared.HasRemoteImage

func (flags *podmanPTFFlags) checkParameters(authFile string) error {
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

	flags.Image.Name, err = utils.ComputePTF(flags.Image.Registry.Host, flags.CustomerID, projectID, serverImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The computed image is %s"), flags.Image.Name)

	images := map[string]*string{
		podman_shared.ServerAttestationService + "@": &flags.Coco.Image.Name,
		podman_shared.HubXmlrpcService:               &flags.HubXmlrpc.Image.Name,
		podman_shared.SalineService + "@":            &flags.Saline.Image.Name,
		podman_shared.DBService:                      &flags.Pgsql.Image.Name,
	}

	for service, pointer := range images {
		if containerImage := getServiceImage(service); containerImage != "" {
			// If no image was found then skip it during the upgrade.
			containerImage, err =
				utils.ComputePTF(flags.Image.Registry.Host, flags.CustomerID, projectID, containerImage, suffix)
			if err != nil {
				return err
			}
			if hasRemoteImage(containerImage, authFile) {
				*pointer = containerImage
				log.Info().Msgf(L("The %[1]s service image is %[2]s"), service, *pointer)
			}
		}
	}

	return nil
}
