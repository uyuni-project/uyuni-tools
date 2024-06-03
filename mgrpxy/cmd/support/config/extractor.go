// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func extract(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	// Copy the generated file locally
	tmpDir, err := os.MkdirTemp("", "mgrpxy-*")
	if err != nil {
		return utils.Errorf(err, L("failed to create temporary directory"))
	}
	defer os.RemoveAll(tmpDir)
	var fileList []string
	if podman.HasService(podman.ServerService) {
		fileList, err = podman.RunSupportConfigOnHost(tmpDir)
	}

	if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		fileList, err = kubernetes.RunSupportConfigOnHost(tmpDir)
	}

	if err != nil {
		return err
	}

	if err := utils.CreateSupportConfigTarball(flags.Output, fileList); err != nil {
		return err
	}

	return nil
}
