// SPDX-FileCopyrightText: 2026 SUSE LLC
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
    "github.com/uyuni-project/uyuni-tools/shared/ssl"
    "github.com/uyuni-project/uyuni-tools/shared/types"
    "github.com/uyuni-project/uyuni-tools/shared/utils"
)

var systemd podman.Systemd = podman.NewSystemd()

func filesRemover(files []string) {
    for _, file := range files {
        if !utils.FileExists(file) {
            log.Trace().Msgf("%s will not removed since it doesn't exists", file)
            continue
        }
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

    var fileList []string
    var supportConfigErr error

    // Run supportconfig but continue even if it fails
    fileList, supportConfigErr = cnx.RunSupportConfig(tmpDir)
    if supportConfigErr != nil {
        log.Warn().Err(supportConfigErr).Msg(L("supportconfig failed, continuing with SSL collection"))
    }

    // Collect SSL certificate information
    sslInfoFile, sslErr := ssl.CollectSSLCertInfo(tmpDir, cnx.Exec)
    if sslErr != nil {
        log.Warn().Err(sslErr).Msg(L("failed to collect SSL certificate information"))
    } else {
        fileList = append(fileList, sslInfoFile)
    }

    var fileListHost []string
    if systemd.HasService(podman.ServerService) {
        fileListHost, err = podman.RunSupportConfigOnPodmanHost(systemd, tmpDir)
    }
    defer filesRemover(fileListHost)
    if err != nil {
        log.Warn().Err(err).Msg(L("failed to run supportconfig on host"))
    }

    if len(fileListHost) > 0 {
        fileList = append(fileList, fileListHost...)
    }

    // If no files collected at all, return the original error
    if len(fileList) == 0 {
        if supportConfigErr != nil {
            return supportConfigErr
        }
        return utils.Errorf(nil, "no files collected")
    }

    return utils.CreateSupportConfigTarball(flags.Output, fileList)
}
