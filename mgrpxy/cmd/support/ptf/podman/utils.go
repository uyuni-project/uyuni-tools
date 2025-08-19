// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
//go:build ptf

package podman

import (
	"errors"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
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
	if _, err := exec.LookPath("podman"); err != nil {
		return errors.New(L("install podman before running this command"))
	}

	if err := updateParameters(flags); err != nil {
		return err
	}

	if err := systemd.StopService(podman_shared.ProxyService); err != nil {
		return err
	}

	hostData, err := podman_shared.InspectHost()
	if err != nil {
		return err
	}

	authFile, cleaner, err := podman_shared.PodmanLogin(hostData, flags.UpgradeFlags.SCC)
	if err != nil {
		return utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	httpdImage, err := podman_shared.PrepareImage(authFile, flags.UpgradeFlags.Httpd.Name,
		flags.UpgradeFlags.PullPolicy, true)
	if err != nil {
		log.Warn().Msgf(L("cannot find httpd image: it will no be upgraded"))
	}
	saltBrokerImage, err := podman_shared.PrepareImage(authFile, flags.UpgradeFlags.SaltBroker.Name,
		flags.UpgradeFlags.PullPolicy, true)
	if err != nil {
		log.Warn().Msgf(L("cannot find salt-broker image: it will no be upgraded"))
	}
	squidImage, err := podman_shared.PrepareImage(authFile, flags.UpgradeFlags.Squid.Name,
		flags.UpgradeFlags.PullPolicy, true)
	if err != nil {
		log.Warn().Msgf(L("cannot find squid image: it will no be upgraded"))
	}
	sshImage, err := podman_shared.PrepareImage(authFile, flags.UpgradeFlags.SSH.Name,
		flags.UpgradeFlags.PullPolicy, true)
	if err != nil {
		log.Warn().Msgf(L("cannot find ssh image: it will no be upgraded"))
	}
	tftpdImage, err := podman_shared.PrepareImage(authFile, flags.UpgradeFlags.Tftpd.Name,
		flags.UpgradeFlags.PullPolicy, true)
	if err != nil {
		log.Warn().Msgf(L("cannot find tftpd image: it will no be upgraded"))
	}

	// Setup the systemd service configuration options
	err = podman.GenerateSystemdService(systemd, httpdImage, saltBrokerImage, squidImage, sshImage,
		tftpdImage, &flags.UpgradeFlags)
	if err != nil {
		return err
	}

	return podman.StartPod(systemd)
}

func updateParameters(flags *podmanPTFFlags) error {
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
		config.imageFlag.Name, err = utils.ComputePTF(flags.UpgradeFlags.ProxyImageFlags.Registry,
			flags.CustomerID, projectID, runningImage, suffix)
		if err != nil {
			return err
		}
		log.Info().Msgf(L("The %[1]s ptf image computed is: %[2]s"), config.serviceName, config.imageFlag.Name)
	}

	return nil
}
