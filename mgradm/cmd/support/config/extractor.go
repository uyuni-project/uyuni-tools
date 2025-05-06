// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func filesRemover(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			log.Error().Err(err).Msgf(L("failed to remove %s temporary file"), file)
		}
	}
}

func extract(_ *types.GlobalFlags, flags *configFlags, _ *cobra.Command, _ []string) error {
	containerName, err := shared.ChooseObjPodmanOrKubernetes(systemd, podman.ServerContainerName, kubernetes.ServerApp)
	if err != nil {
		return err
	}

	cnx := shared.NewConnection(flags.Backend, containerName, kubernetes.ServerFilter)

	// Copy the generated file locally
	tmpDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	fileList, err := cnx.RunSupportConfig(tmpDir)
	if err != nil {
		return err
	}

	var fileListHost []string
	if systemd.HasService(podman.ServerService) {
		fileListHost, err = podman.RunSupportConfigOnPodmanHost(systemd, tmpDir)
	}
	defer filesRemover(fileListHost)
	if err != nil {
		return err
	}

	if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		var namespace string
		namespace, err = cnx.GetNamespace("")
		if err != nil {
			return err
		}
		fileListHost, err = kubernetes.RunSupportConfigOnKubernetesHost(tmpDir, namespace, kubernetes.ServerFilter)
	}
	if err != nil {
		return err
	}

	if len(fileListHost) > 0 {
		fileList = append(fileList, fileListHost...)
	}

	return utils.CreateSupportConfigTarball(flags.Output, fileList)
}
