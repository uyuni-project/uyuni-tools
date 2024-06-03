// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"errors"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
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
	if err := flags.checkParameters(); err != nil {
		return err
	}
	return podman.Upgrade(systemd, globalFlags, &flags.UpgradeFlags, cmd, args)
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
	suffix := "ptf"
	projectId := flags.PTFId
	if flags.TestId != "" {
		suffix = "test"
		projectId = flags.TestId
	}
	httpdImage, err := podman_shared.GetRunningImage("httpd")
	if err != nil {
		return err
	}
	flags.UpgradeFlags.Httpd.Name, err = utils.ComputePTF(flags.CustomerId, projectId, httpdImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The httpd ptf image computed is: %s"), flags.UpgradeFlags.Httpd.Name)

	sshImage, err := podman_shared.GetRunningImage("ssh")
	if err != nil {
		return err
	}
	flags.UpgradeFlags.Ssh.Name, err = utils.ComputePTF(flags.CustomerId, projectId, sshImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The ssh ptf image computed is: %s"), flags.UpgradeFlags.Ssh.Name)

	tftpdImage, err := podman_shared.GetRunningImage("tftpd")
	if err != nil {
		return err
	}
	flags.UpgradeFlags.Tftpd.Name, err = utils.ComputePTF(flags.CustomerId, projectId, tftpdImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The tftpd ptf image computed is: %s"), flags.UpgradeFlags.Tftpd.Name)

	saltBrokerImage, err := podman_shared.GetRunningImage("salt-broker")
	if err != nil {
		return err
	}
	flags.UpgradeFlags.SaltBroker.Name, err = utils.ComputePTF(flags.CustomerId, projectId, saltBrokerImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The salt-broker ptf image computed is: %s"), flags.UpgradeFlags.SaltBroker.Name)

	squidImage, err := podman_shared.GetRunningImage("squid")
	if err != nil {
		return err
	}
	flags.UpgradeFlags.Squid.Name, err = utils.ComputePTF(flags.CustomerId, projectId, squidImage, suffix)
	if err != nil {
		return err
	}
	log.Info().Msgf(L("The squid ptf image computed is: %s"), flags.UpgradeFlags.Squid.Name)

	return nil
}
