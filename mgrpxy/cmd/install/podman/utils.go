// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Start the proxy services.
func startPod() error {
	ret := shared_podman.IsServiceRunning(shared_podman.ProxyService)
	if ret {
		return shared_podman.RestartService(shared_podman.ProxyService)
	} else {
		return shared_podman.EnableService(shared_podman.ProxyService)
	}
}

func installForPodman(globalFlags *types.GlobalFlags, flags *podmanProxyInstallFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf("install podman before running this command")
	}

	configPath := utils.GetConfigPath(args)
	if err := unpackConfig(configPath); err != nil {
		return fmt.Errorf("failed to extract proxy config from %s file: %s", configPath, err)
	}

	httpdImage, err := getContainerImage(flags, "httpd")
	if err != nil {
		return err
	}
	saltBrokerImage, err := getContainerImage(flags, "salt-broker")
	if err != nil {
		return err
	}
	squidImage, err := getContainerImage(flags, "squid")
	if err != nil {
		return err
	}
	sshImage, err := getContainerImage(flags, "ssh")
	if err != nil {
		return err
	}
	tftpdImage, err := getContainerImage(flags, "tftpd")
	if err != nil {
		return err
	}

	// Setup the systemd service configuration options
	if err := podman.GenerateSystemdService(httpdImage, saltBrokerImage, squidImage, sshImage, tftpdImage, flags.Podman.Args); err != nil {
		return fmt.Errorf("cannot generate systemd file: %s", err)
	}

	return startPod()
}

func getContainerImage(flags *podmanProxyInstallFlags, name string) (string, error) {
	image := flags.GetContainerImage(name)
	inspectedHostValues, err := adm_utils.InspectHost()
	if err != nil {
		return "", fmt.Errorf("cannot inspect host values: %s", err)
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	if err := shared_podman.PrepareImage(image, flags.PullPolicy, pullArgs...); err != nil {
		return "", err
	}
	return image, nil
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
