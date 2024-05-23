// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

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
	Args []string `mapstructure:"arg"`
}

// GetCommonParams splits the common arguments.
func GetCommonParams() []string {
	return strings.Split(commonArgs, " ")
}

// AddPodmanArgFlag add the podman arguments to a command.
func AddPodmanArgFlag(cmd *cobra.Command) {
	cmd.Flags().StringSlice("podman-arg", []string{}, L("Extra arguments to pass to podman"))
}

// EnablePodmanSocket enables the podman socket.
func EnablePodmanSocket() error {
	err := utils.RunCmd("systemctl", "enable", "--now", "podman.socket")
	if err != nil {
		return utils.Errorf(err, L("failed to enable podman.socket unit"))
	}
	return err
}

// RunContainer execute a container.
func RunContainer(name string, image string, extraArgs []string, cmd []string) error {
	podmanArgs := append([]string{"run", "--name", name}, GetCommonParams()...)
	podmanArgs = append(podmanArgs, extraArgs...)
	for _, volume := range utils.ServerVolumeMounts {
		podmanArgs = append(podmanArgs, "-v", volume.Name+":"+volume.MountPath)
	}
	podmanArgs = append(podmanArgs, image)
	podmanArgs = append(podmanArgs, cmd...)

	err := utils.RunCmdStdMapping(zerolog.DebugLevel, "podman", podmanArgs...)
	if err != nil {
		return utils.Errorf(err, L("failed to run %s container"), name)
	}

	return nil
}

// DeleteContainer deletes a container based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteContainer(name string, dryRun bool) {
	if out, _ := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "ps", "-a", "-q", "-f", "name="+name); len(out) > 0 {
		if dryRun {
			log.Info().Msgf(L("Would run podman kill %[1]s for container id %[2]s"), name, out)
			log.Info().Msgf(L("Would run podman remove %[1]s for container id %[2]s"), name, out)
		} else {
			log.Info().Msgf(L("Run podman kill %[1]s for container id %[2]s"), name, out)
			err := utils.RunCmd("podman", "kill", name)
			if err != nil {
				log.Error().Err(err).Msg(L("Failed to kill the server"))

				log.Info().Msgf(L("Run podman remove %[1]s for container id %[2]s"), name, out)
				err = utils.RunCmd("podman", "rm", name)
				if err != nil {
					log.Error().Err(err).Msg(L("Error removing container"))
				}
			}
		}
	} else {
		log.Info().Msg(L("Container already removed"))
	}
}

// DeleteVolume deletes a podman volume based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteVolume(name string, dryRun bool) error {
	exists := isVolumePresent(name)
	if exists {
		if dryRun {
			log.Info().Msgf(L("Would run %s"), "podman volume rm "+name)
		} else {
			log.Info().Msgf(L("Run %s"), "podman volume rm "+name)
			err := utils.RunCmd("podman", "volume", "rm", name)
			if err != nil {
				log.Error().Err(err).Msgf(L("Failed to remove volume %s"), name)
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

// Inspect check values on a given image and deploy.
func Inspect(serverImage string, pullPolicy string, proxyHost bool) (map[string]string, error) {
	scriptDir, err := os.MkdirTemp("", "mgradm-*")
	defer os.RemoveAll(scriptDir)
	if err != nil {
		return map[string]string{}, utils.Errorf(err, L("failed to create temporary directory"))
	}

	inspectedHostValues, err := utils.InspectHost(proxyHost)
	if err != nil {
		return map[string]string{}, utils.Errorf(err, L("cannot inspect host values"))
	}

	pullArgs := []string{}
	_, scc_user_exist := inspectedHostValues["host_scc_username"]
	_, scc_user_password := inspectedHostValues["host_scc_password"]
	if scc_user_exist && scc_user_password {
		pullArgs = append(pullArgs, "--creds", inspectedHostValues["host_scc_username"]+":"+inspectedHostValues["host_scc_password"])
	}

	preparedImage, err := PrepareImage(serverImage, pullPolicy, pullArgs...)
	if err != nil {
		return map[string]string{}, err
	}

	if err := utils.GenerateInspectContainerScript(scriptDir); err != nil {
		return map[string]string{}, err
	}

	podmanArgs := []string{
		"-v", scriptDir + ":" + utils.InspectOutputFile.Directory,
		"--security-opt", "label:disable",
	}

	err = RunContainer("uyuni-inspect", preparedImage, podmanArgs,
		[]string{utils.InspectOutputFile.Directory + "/" + utils.InspectScriptFilename})
	if err != nil {
		return map[string]string{}, err
	}

	inspectResult, err := utils.ReadInspectData(scriptDir)
	if err != nil {
		return map[string]string{}, utils.Errorf(err, L("cannot inspect data"))
	}

	return inspectResult, err
}
