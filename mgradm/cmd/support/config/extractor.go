// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func extract(globalFlags *types.GlobalFlags, flags *configFlags, cmd *cobra.Command, args []string) error {
	cnx := shared.NewConnection(flags.Backend, podman.ServerContainerName, kubernetes.ServerFilter)

	// Copy the generated file locally
	tmpDir, err := os.MkdirTemp("", "mgradm-*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	var files []string
	extensions := []string{"", ".md5"}

	// Run supportconfig in the container if it's running
	log.Info().Msg("Running supportconfig in the container")
	out, err := cnx.Exec("supportconfig")
	if err != nil {
		log.Err(err).Msg("Failed to run supportconfig")
	} else {
		tarballPath := getSupportConfigPath(out)
		if tarballPath == "" {
			return fmt.Errorf("failed to find container supportconfig tarball from command output")
		}

		// TODO Get the error from copy
		for _, ext := range extensions {
			containerTarball := path.Join(tmpDir, "container-supportconfig.txz"+ext)
			cnx.Copy("server:"+tarballPath+ext, containerTarball, "", "")
			files = append(files, containerTarball)

			// Remove the generated file in the container
			if _, err := cnx.Exec("rm", tarballPath+ext); err != nil {
				log.Warn().Err(err).Msgf("Failed to remove %s%s file in the container", tarballPath, ext)
			}
		}
	}

	// Run supportconfig on the host if installed
	if _, err := exec.LookPath("supportconfig"); err == nil {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "supportconfig")
		if err != nil {
			log.Warn().Err(err).Msg("Failed to run supportconfig on the host")
		}
		tarballPath := getSupportConfigPath(out)

		// Look for the generated supportconfig file
		if tarballPath != "" && utils.FileExists(tarballPath) {
			for _, ext := range extensions {
				files = append(files, tarballPath+ext)
			}
		} else {
			log.Error().Msg("Failed to find host supportconfig tarball from command output")
		}
	} else {
		log.Warn().Msg("supportconfig is not available on the host, skipping it")
	}

	// TODO Get cluster infos in case of kubernetes

	// Pack it all into a tarball
	log.Info().Msg("Preparing the tarball")
	tarball, err := utils.NewTarGz(flags.Output)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := tarball.AddFile(file, path.Base(file)); err != nil {
			return fmt.Errorf("failed to add %s to tarball: %s", path.Base(file), err)
		}
	}
	tarball.Close()

	return nil
}

func getSupportConfigPath(out []byte) string {
	re := regexp.MustCompile(`/var/log/scc_[^.]+\.txz`)
	return re.FindString(string(out))
}
