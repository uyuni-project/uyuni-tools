// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func extract(_ *types.GlobalFlags, flags *configFlags, _ *cobra.Command, _ []string) error {
	// Copy the generated file locally
	tmpDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	var fileList []string
	if podman.HasService(podman.ProxyService) {
		fileList, err = podman.RunSupportConfigOnPodmanHost(tmpDir)
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
