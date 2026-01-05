// SPDX-FileCopyrightText: 2026 SUSE LLC
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

	hostData, err := podman_shared.InspectHost()
	if err != nil {
		return err
	}
	authFile, cleaner, err := podman_shared.PodmanLogin(hostData, flags.UpgradeFlags.Registry, flags.SCC)
	if err != nil {
		return err
	}
	defer cleaner()

	if err := updateParameters(flags, authFile); err != nil {
		return err
	}

	if err := systemd.StopService(podman_shared.ProxyService); err != nil {
		return err
	}

	pullPolicy := flags.UpgradeFlags.PullPolicy
	httpdImage := getImage(authFile, flags.UpgradeFlags.Httpd.Name, pullPolicy)
	saltBrokerImage := getImage(authFile, flags.UpgradeFlags.SaltBroker.Name, pullPolicy)
	squidImage := getImage(authFile, flags.UpgradeFlags.Squid.Name, pullPolicy)
	sshImage := getImage(authFile, flags.UpgradeFlags.SSH.Name, pullPolicy)
	tftpdImage := getImage(authFile, flags.UpgradeFlags.Tftpd.Name, pullPolicy)

	// Setup the systemd service configuration options
	err = podman.GenerateSystemdService(systemd, httpdImage, saltBrokerImage, squidImage, sshImage,
		tftpdImage, &flags.UpgradeFlags)
	if err != nil {
		return err
	}

	return podman.StartPod(systemd)
}

func getImage(authFile string, image string, policy string) string {
	var newImage string
	var err error
	if image != "" {
		newImage, err = podman_shared.PrepareImage(authFile, image, policy, true)
		if err != nil {
			log.Warn().Msgf(L("cannot find %s image: it will no be upgraded"), image)
		}
	}

	return newImage
}

// variables for unit testing.
var getServiceImage = podman_shared.GetServiceImage
var hasRemoteImage = podman_shared.HasRemoteImage

func checkIDs(flags *podmanPTFFlags) error {
	if flags.TestID != "" && flags.PTFId != "" {
		return errors.New(L("ptf and test flags cannot be set simultaneously "))
	}
	if flags.TestID == "" && flags.PTFId == "" {
		return errors.New(L("ptf and test flags cannot be empty simultaneously "))
	}
	if flags.CustomerID == "" {
		return errors.New(L("user flag cannot be empty"))
	}
	return nil
}

func updateParameters(flags *podmanPTFFlags, authFile string) error {
	if err := checkIDs(flags); err != nil {
		return err
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
		{podman.ServiceHTTPd, &flags.UpgradeFlags.Httpd},
		{podman.ServiceSSH, &flags.UpgradeFlags.SSH},
		{podman.ServiceTFTFd, &flags.UpgradeFlags.Tftpd},
		{podman.ServiceSaltBroker, &flags.UpgradeFlags.SaltBroker},
		{podman.ServiceSquid, &flags.UpgradeFlags.Squid},
	}
	// Process each pxy image
	for _, config := range proxyImages {
		if containerImage := getServiceImage(config.serviceName); containerImage != "" {
			// If no image was found then skip it during the upgrade.
			newImage, err := utils.ComputePTF(flags.UpgradeFlags.ProxyImageFlags.Registry.Host, flags.CustomerID, projectID,
				containerImage, suffix)
			log.Debug().Msgf("computed PTF image url: %s", newImage)
			if err != nil {
				return err
			}
			if hasRemoteImage(newImage, authFile) {
				config.imageFlag.Name = newImage
				log.Info().Msgf(L("The %[1]s service image is %[2]s"), config.serviceName, newImage)
			} else {
				config.imageFlag.Name = ""
			}
		}
	}

	return nil
}
