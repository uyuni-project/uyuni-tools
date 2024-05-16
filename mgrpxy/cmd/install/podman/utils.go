// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/podman"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
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

func installForPodman(globalFlags *types.GlobalFlags, flags *podman.PodmanProxyFlags, cmd *cobra.Command, args []string) error {
	if _, err := exec.LookPath("podman"); err != nil {
		return fmt.Errorf(L("install podman before running this command"))
	}

	configPath := utils.GetConfigPath(args)
	if err := podman.UnpackConfig(configPath); err != nil {
		return shared_utils.Errorf(err, L("failed to extract proxy config from %s file"), configPath)
	}

	httpdImage, err := podman.GetContainerImage(&flags.ProxyImageFlags, "httpd")
	if err != nil {
		return err
	}
	saltBrokerImage, err := podman.GetContainerImage(&flags.ProxyImageFlags, "salt-broker")
	if err != nil {
		return err
	}
	squidImage, err := podman.GetContainerImage(&flags.ProxyImageFlags, "squid")
	if err != nil {
		return err
	}
	sshImage, err := podman.GetContainerImage(&flags.ProxyImageFlags, "ssh")
	if err != nil {
		return err
	}
	tftpdImage, err := podman.GetContainerImage(&flags.ProxyImageFlags, "tftpd")
	if err != nil {
		return err
	}

	// Setup the systemd service configuration options
	if err := podman.GenerateSystemdService(httpdImage, saltBrokerImage, squidImage, sshImage, tftpdImage, flags); err != nil {
		return err
	}

	return startPod()
}
