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

func ptfForPodman(
	globalFlags *types.GlobalFlags,
	flags *podmanPTFFlags,
	cmd *cobra.Command,
	args []string,
) error {
	//we don't want to perform a postgres version upgrade when installing a PTF.
	//in that case, we can use the upgrade command.
	dummyMigration := types.ImageFlags{}
	dummyCoco := types.ImageFlags{}
	if err := flags.checkParameters(); err != nil {
		return err
	}

	authFile, cleaner, err := podman_shared.PodmanLogin()
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	return podman.Upgrade(authFile, "", flags.Image, dummyMigration, dummyCoco)
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
	if flags.TestId != "" {
		suffix = "test"
	}
	flags.Image.Name, err = utils.ComputePTF(flags.CustomerId, flags.PTFId, serverImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The image computed is: %s"), flags.Image.Name)
	return nil
}
