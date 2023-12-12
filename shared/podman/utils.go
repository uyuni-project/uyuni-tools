// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const ServerContainerName = "uyuni-server"

var ProxyContainerNames = []string{
	"uyuni-proxy-httpd",
	"uyuni-proxy-salt-broker",
	"uyuni-proxy-squid",
	"uyuni-proxy-ssh",
	"uyuni-proxy-tftpd",
}

type PodmanFlags struct {
	Args []string `mapstructure:"arg"`
}

func AddPodmanInstallFlag(cmd *cobra.Command) {
	cmd.Flags().StringSlice("podman-arg", []string{}, "Extra arguments to pass to podman")
}

func EnablePodmanSocket() {
	err := utils.RunCmd("systemctl", "enable", "--now", "podman.socket")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to enable podman.socket unit")
	}
}

// DeleteContainer deletes a container based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteContainer(name string, dryRun bool) {
	if out, _ := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "ps", "-a", "-q", "-f", "name="+name); len(out) > 0 {
		if dryRun {
			log.Info().Msgf("Would run podman kill %s for container id: %s", name, out)
			log.Info().Msgf("Would run podman remove %s for container id: %s", name, out)
		} else {
			log.Info().Msgf("Run podman kill %s for container id: %s", name, out)
			err := utils.RunCmd("podman", "kill", name)
			if err != nil {
				log.Error().Err(err).Msg("Failed to kill the server")

				log.Info().Msgf("Run podman remove %s for container id: %s", name, out)
				err = utils.RunCmd("podman", "rm", name)
				if err != nil {
					log.Error().Err(err).Msg("Error removing container")
				}
			}
		}
	} else {
		log.Info().Msg("Container already removed")
	}
}

// DeleteVolume deletes a podman volume based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteVolume(name string, dryRun bool) {
	cmd := exec.Command("podman", "volume", "exists", name)
	cmd.Run()
	if cmd.ProcessState.ExitCode() == 0 {
		if dryRun {
			log.Info().Msgf("Would run podman volume rm %s", name)
		} else {
			log.Info().Msgf("Run podman volume rm %s", name)
			err := utils.RunCmd("podman", "volume", "rm", name)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to remove volume %s", name)
			}
		}
	}
}
