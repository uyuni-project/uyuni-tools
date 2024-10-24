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
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
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
	// we don't want to perform a postgres version upgrade when installing a PTF.
	// in that case, we can use the upgrade command.
	dummyImage := types.ImageFlags{}
	dummyCoco := adm_utils.CocoFlags{}
	dummyHubXmlrpc := adm_utils.HubXmlrpcFlags{}
	if err := flags.checkParameters(); err != nil {
		return err
	}

	hostData, err := podman_shared.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman_shared.PodmanLogin(hostData, flags.SCC)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	return podman.Upgrade(systemd, authFile, "", flags.Image, dummyImage, dummyCoco, dummyHubXmlrpc)
}

func (flags *podmanPTFFlags) checkParameters() error {
	if flags.TestId != "" && flags.PTFId != "" {
		return errors.New(L("ptf and test flags cannot be set simultaneously "))
	}
	if flags.TestId == "" && flags.PTFId == "" {
		return errors.New(L("ptf and test flags cannot be empty simultaneously "))
	}
	if flags.CustomerId == "" {
		return errors.New(L("user flag cannot be empty"))
	}
	serverImage, err := podman_shared.GetRunningImage(podman_shared.ServerContainerName)
	if err != nil {
		return err
	}

	suffix := "ptf"
	projectId := flags.PTFId
	if flags.TestId != "" {
		suffix = "test"
		projectId = flags.TestId
	}
	flags.Image.Name, err = utils.ComputePTF(flags.CustomerId, projectId, serverImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The image computed is: %s"), flags.Image.Name)
	return nil
}
