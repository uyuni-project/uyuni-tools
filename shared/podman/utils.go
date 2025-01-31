// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// runCmdOutput is a function pointer to use for easies unit testing.
var runCmdOutput = utils.RunCmdOutput

const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

// ServerContainerName represents the server container name.
const ServerContainerName = "uyuni-server"

// HubXmlrpcContainerName is the container name for the Hub XML-RPC API.
const HubXmlrpcContainerName = "uyuni-hub-xmlrpc"

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
func RunContainer(name string, image string, volumes []types.VolumeMount, extraArgs []string, cmd []string) error {
	podmanArgs := append([]string{"run", "--name", name}, GetCommonParams()...)
	podmanArgs = append(podmanArgs, extraArgs...)
	for _, volume := range volumes {
		podmanArgs = append(podmanArgs, "-v", volume.Name+":"+volume.MountPath)
	}
	podmanArgs = append(podmanArgs, "--network", UyuniNetwork)
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

// GetServiceImage returns the value of the UYUNI_IMAGE variable for a systemd service.
func GetServiceImage(service string) string {
	out, err := runCmdOutput(zerolog.DebugLevel, "systemctl", "cat", service)
	if err != nil {
		log.Warn().Err(err).Msgf(L("failed to get %s systemd service definition"), service)
		return ""
	}

	imageFinder := regexp.MustCompile(`UYUNI_IMAGE=(.*)`)
	matches := imageFinder.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		log.Warn().Msgf(L("no UYUNI_IMAGE defined in %s systemd service"), service)
		return ""
	}
	return matches[1]
}

// DeleteImage deletes a podman image based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteImage(name string, dryRun bool) error {
	exists := imageExists(name)
	if exists {
		if dryRun {
			log.Info().Msgf(L("Would run %s"), "podman image rm "+name)
		} else {
			log.Info().Msgf(L("Run %s"), "podman image rm "+name)
			err := utils.RunCmd("podman", "image", "rm", name)
			if err != nil {
				log.Error().Err(err).Msgf(L("Failed to remove image %s"), name)
			}
		}
	}
	return nil
}

func imageExists(volume string) bool {
	cmd := exec.Command("podman", "image", "exists", volume)
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
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
			if err := utils.RunCmd("podman", "volume", "rm", name); err != nil {
				log.Trace().Err(err).Msgf("podman volume rm %s", name)
				// Check if the volume is not mounted - for example var-pgsql - as second storage device
				// We need to compute volume path ourselves because above `podman volume rm` call may have
				// already removed volume from podman internal structures
				basePath, errBasePath := getPodmanVolumeBasePath()
				if errBasePath != nil {
					return errBasePath
				}
				target := path.Join(basePath, name)
				if isVolumePathEmpty(target) && isVolumePathMounted(target) {
					log.Info().Msgf(L("Volume %s is externally mounted, directory cannot be removed"), name)
					return nil
				}
				return err
			}
		}
	}
	return nil
}

func isVolumePresent(volume string) bool {
	var exitError *exec.ExitError
	cmd := exec.Command("podman", "volume", "exists", volume)
	if err := cmd.Run(); err != nil && errors.As(err, &exitError) {
		log.Debug().Err(err).Msgf("podman volume exists %s", volume)
		return false
	}
	return cmd.ProcessState.Success()
}

func isVolumePathMounted(volume string) bool {
	cmd := exec.Command("findmnt", "--target", volume)
	var exitError *exec.ExitError
	if err := cmd.Run(); err != nil && errors.As(err, &exitError) {
		log.Debug().Err(err).Msgf("findmnt --target %s", volume)
		return false
	}
	return cmd.ProcessState.Success()
}

func isVolumePathEmpty(volume string) bool {
	f, err := os.Open(volume)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	return errors.Is(err, io.EOF)
}

func getPodmanVolumeBasePath() (string, error) {
	cmd := exec.Command("podman", "system", "info", "--format={{ .Store.VolumePath }}")
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// Inspect check values on a given image and deploy.
func Inspect(
	serverImage string,
	pullPolicy string,
	scc types.SCCCredentials,
) (*utils.ServerInspectData, error) {
	scriptDir, cleaner, err := utils.TempDir()
	if err != nil {
		return nil, err
	}
	defer cleaner()

	hostData, err := InspectHost()
	if err != nil {
		return nil, err
	}

	authFile, cleaner, err := PodmanLogin(hostData, scc)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to login to registry.suse.com"))
	}
	defer cleaner()

	preparedImage, err := PrepareImage(authFile, serverImage, pullPolicy, true)
	if err != nil {
		return nil, err
	}

	inspector := utils.NewServerInspector(scriptDir)
	if err := inspector.GenerateScript(); err != nil {
		return nil, err
	}

	podmanArgs := []string{
		"-v", scriptDir + ":" + utils.InspectContainerDirectory,
		"--security-opt", "label=disable",
	}

	err = RunContainer("uyuni-inspect", preparedImage, utils.ServerVolumeMounts, podmanArgs,
		[]string{utils.InspectContainerDirectory + "/" + utils.InspectScriptFilename})
	if err != nil {
		return nil, err
	}

	inspectResult, err := inspector.ReadInspectData()
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect data"))
	}

	return inspectResult, err
}
