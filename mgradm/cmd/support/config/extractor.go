// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func extract(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)

	// Copy the generated file locally
	tmpDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return utils.Errorf(err, L("failed to create temporary directory"))
	}
	defer os.RemoveAll(tmpDir)

	var files []string
	extensions := []string{"", ".md5"}

	// Run supportconfig in the container if it's running
	log.Info().Msg(L("Running supportconfig in the container"))
	out, err := cnx.Exec("supportconfig")
	if err != nil {
		return errors.New(L("failed to run supportconfig"))
	} else {
		tarballPath := utils.GetSupportConfigPath(out)
		if tarballPath == "" {
			return fmt.Errorf(L("failed to find container supportconfig tarball from command output"))
		}

		// TODO Get the error from copy
		for _, ext := range extensions {
			containerTarball := path.Join(tmpDir, "container-supportconfig.txz"+ext)
			if err := cnx.Copy("server:"+tarballPath+ext, containerTarball, "", ""); err != nil {
				return utils.Errorf(err, L("cannot copy tarball"))
			}
			files = append(files, containerTarball)

			// Remove the generated file in the container
			if _, err := cnx.Exec("rm", tarballPath+ext); err != nil {
				return utils.Errorf(err, L("failed to remove %s file in the container"), tarballPath+ext)
			}
		}
	}

	// Run supportconfig on the host if installed
	if _, err := exec.LookPath("supportconfig"); err == nil {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "supportconfig")
		if err != nil {
			return utils.Errorf(err, L("failed to run supportconfig on the host"))
		}
		tarballPath := utils.GetSupportConfigPath(out)

		// Look for the generated supportconfig file
		if tarballPath != "" && utils.FileExists(tarballPath) {
			for _, ext := range extensions {
				files = append(files, tarballPath+ext)
			}
		} else {
			return errors.New(L("failed to find host supportconfig tarball from command output"))
		}
	} else {
		log.Warn().Msg(L("supportconfig is not available on the host, skipping it"))
	}

	// TODO Get cluster infos in case of kubernetes

	// Pack it all into a tarball
	log.Info().Msg(L("Preparing the tarball"))

	supportFileName := utils.GetSupportConfigFileSaveName()
	supportFilePath := path.Join(flags.Output, fmt.Sprintf("%s.tar.gz", supportFileName))

	tarball, err := utils.NewTarGz(supportFilePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := tarball.AddFile(file, path.Join(supportFileName, path.Base(file))); err != nil {
			return utils.Errorf(err, L("failed to add %s to tarball"), path.Base(file))
		}
	}
	tarball.Close()

	return nil
}
