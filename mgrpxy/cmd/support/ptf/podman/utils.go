// SPDX-FileCopyrightText: 2025 SUSE LLC
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

	proxyImages := []struct {
		serviceName string
		imageFlag   *types.ImageFlags
	}{
		{"httpd", &flags.UpgradeFlags.Httpd},
		{"ssh", &flags.UpgradeFlags.SSH},
		{"tftpd", &flags.UpgradeFlags.Tftpd},
		{"salt-broker", &flags.UpgradeFlags.SaltBroker},
		{"squid", &flags.UpgradeFlags.Squid},
	}

	// Process each pxy image
	for _, config := range proxyImages {
		runningImage, err := podman_shared.GetRunningImage(config.serviceName)
		if err != nil {
			return err
		}

		config.imageFlag.Name, err = utils.ComputePTF(flags.UpgradeFlags.SCC.Registry, flags.CustomerID, projectID,
			runningImage, suffix)
		if err != nil {
			return err
		}
		config.imageFlag.SkipComputation = true

		log.Info().Msgf(L("The %[1]s ptf image computed is: %[2]s"), config.serviceName, config.imageFlag.Name)
	}

	return nil
}
