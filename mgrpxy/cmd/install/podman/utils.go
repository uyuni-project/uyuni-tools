// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Start the proxy services.
func startPod() {
	const servicePod = "uyuni-proxy-pod"
	if shared_podman.IsServiceRunning(servicePod) {
		shared_podman.RestartService(servicePod)
	} else {
		shared_podman.EnableService(servicePod)
	}
}

func installForPodman(globalFlags *types.GlobalFlags, flags *podmanProxyInstallFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		log.Fatal().Err(err).Msgf("install podman before running this command")
	}

	configPath := utils.GetConfigPath(args)
	if err := unpackConfig(configPath); err != nil {
		log.Fatal().Err(err).Msgf("Failed to extract proxy config from %s file", configPath)
	}

	httpdImage := getContainerImage(flags, "httpd")
	saltBrokerImage := getContainerImage(flags, "salt-broker")
	squidImage := getContainerImage(flags, "squid")
	sshImage := getContainerImage(flags, "ssh")
	tftpdImage := getContainerImage(flags, "tftpd")

	// Setup the systemd service configuration options
	podman.GenerateSystemdService(httpdImage, saltBrokerImage, squidImage, sshImage, tftpdImage, flags.Podman.Args)

	startPod()
	return nil
}

func getContainerImage(flags *podmanProxyInstallFlags, name string) string {
	image := flags.GetContainerImage(name)
	shared_podman.PrepareImage(image, flags.PullPolicy)
	return image
}

func unpackConfig(configPath string) error {
	log.Info().Msgf("Setting up proxy with configuration %s", configPath)
	const proxyConfigDir = "/etc/uyuni/proxy"
	if err := os.MkdirAll(proxyConfigDir, 0755); err != nil {
		return err
	}

	if err := shared_utils.ExtractTarGz(configPath, proxyConfigDir); err != nil {
		return err
	}
	return nil
}
