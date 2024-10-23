// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func extract(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	// Copy the generated file locally
	tmpDir, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	var fileList []string
	if systemd.HasService(podman.ProxyService) {
		fileList, err = podman.RunSupportConfigOnPodmanHost(systemd, tmpDir)
	}

	if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		cnx := shared.NewConnection("kubectl", "", kubernetes.ProxyFilter)
		var namespace string
		namespace, err = cnx.GetNamespace("")
		if err != nil {
			return err
		}
		fileList, err = kubernetes.RunSupportConfigOnKubernetesHost(tmpDir, namespace, kubernetes.ProxyFilter)
	}

	if err != nil {
		return err
	}

	if err := utils.CreateSupportConfigTarball(flags.Output, fileList); err != nil {
		return err
	}

	return nil
}
