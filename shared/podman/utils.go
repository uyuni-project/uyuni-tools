// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ServerContainerName represents the server container name.
const ServerContainerName = "uyuni-server"

// ProxyContainerNames represents all the proxy container names.
var ProxyContainerNames = []string{
	"uyuni-proxy-httpd",
	"uyuni-proxy-salt-broker",
	"uyuni-proxy-squid",
	"uyuni-proxy-ssh",
	"uyuni-proxy-tftpd",
}

// PodmanFlags stores the podman arguments.
type PodmanFlags struct {
	Args   []string         `mapstructure:"arg"`
	Mounts PodmanMountFlags `mapstructure:"mount"`
}

// PodmanMountFlags stores the --podman-mount-* arguments.
type PodmanMountFlags struct {
	Cache      string
	Postgresql string
	Spacewalk  string
}

// AddPodmanInstallFlag add the podman arguments to a command.
func AddPodmanInstallFlag(cmd *cobra.Command) {
	cmd.Flags().StringSlice("podman-arg", []string{}, "Extra arguments to pass to podman")

	cmd.Flags().String("podman-mount-cache", "", "Path to custom /var/cache volume")
	cmd.Flags().String("podman-mount-postgresql", "", "Path to custom /var/lib/pgsql volume")
	cmd.Flags().String("podman-mount-spacewalk", "", "Path to custom /var/spacewalk volume")
}

// EnablePodmanSocket enables the podman socket.
func EnablePodmanSocket() error {
	err := utils.RunCmd("systemctl", "enable", "--now", "podman.socket")
	if err != nil {
		return fmt.Errorf("failed to enable podman.socket unit: %s", err)
	}
	return err
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
func DeleteVolume(name string, dryRun bool) error {
	exists := isVolumePresent(name)
	if exists {
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
	return nil
}

func isVolumePresent(volume string) bool {
	cmd := exec.Command("podman", "volume", "exists", volume)
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

// LinkVolumes adds the symlinks for the podman volumes if needed.
func LinkVolumes(mountFlags *PodmanMountFlags) error {
	graphRoot, err := getGraphRoot()
	if err != nil {
		return err
	}

	data := map[string]string{
		"var-cache":     mountFlags.Cache,
		"var-spacewalk": mountFlags.Spacewalk,
		"var-pgsql":     mountFlags.Postgresql,
	}
	for volume, value := range data {
		if value != "" {
			volumePath := path.Join(graphRoot, "volumes", volume)
			if utils.FileExists(volumePath) {
				return fmt.Errorf("volume folder (%s) already exists, cannot link it to %s", volumePath, value)
			}
			baseFolder := path.Join(graphRoot, "volumes")
			if err := utils.RunCmd("mkdir", "-p", baseFolder); err != nil {
				return fmt.Errorf("failed to create volumes folder: %s: %s", baseFolder, err)
			}

			if err := utils.RunCmd("ln", "-s", value, volumePath); err != nil {
				return fmt.Errorf("failed to link volume folder %s to %s: %s", value, volumePath, err)
			}
		}
	}
	return nil
}

func getGraphRoot() (string, error) {
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "system", "info", "--format", "{{ .Store.GraphRoot }}")
	if err != nil {
		return "", fmt.Errorf("failed to get podman's volumes folder: %s", err)
	}
	return strings.TrimSpace(string(out)), nil
}
