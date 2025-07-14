// SPDX-FileCopyrightText: 2025 SUSE LLC
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
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// runCmd* are a function pointers to use for easies unit testing.
var runCmdOutput = utils.RunCmdOutput
var runCmd = utils.RunCmd

const commonArgs = "--rm --cap-add NET_RAW --tmpfs /run -v cgroup:/sys/fs/cgroup:rw"

// ServerContainerName represents the server container name.
const ServerContainerName = "uyuni-server"

// HubXmlrpcContainerName is the container name for the Hub XML-RPC API.
const HubXmlrpcContainerName = "uyuni-hub-xmlrpc"

// DBContainerName represents the database container name.
const DBContainerName = "uyuni-db"

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

// ReadFromContainer read a file from a container.
func ReadFromContainer(name string, image string, volumes []types.VolumeMount,
	extraArgs []string, file string) ([]byte, error) {
	podmanArgs := append([]string{"run", "--name", name}, GetCommonParams()...)
	podmanArgs = append(podmanArgs, extraArgs...)
	for _, volume := range volumes {
		if IsVolumePresent(volume.Name) {
			podmanArgs = append(podmanArgs, "-v", volume.Name+":"+volume.MountPath)
		}
	}
	podmanArgs = append(podmanArgs, "--network", UyuniNetwork)
	podmanArgs = append(podmanArgs, image)
	podmanArgs = append(podmanArgs, []string{"cat", file}...)

	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "podman", podmanArgs...)
	if err != nil {
		return []byte{}, utils.Errorf(err, L("failed to run %s container"), name)
	}

	return out, nil
}

// RunContainer execute a container.
func RunContainer(name string, image string, volumes []types.VolumeMount, extraArgs []string, cmd []string) error {
	podmanArgs := append([]string{"run", "--name", name}, GetCommonParams()...)
	podmanArgs = append(podmanArgs, extraArgs...)
	podmanArgs = append(podmanArgs, "--shm-size=0")
	podmanArgs = append(podmanArgs, "--shm-size-systemd=0")
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

	imageFinder := regexp.MustCompile(`UYUNI.*_IMAGE=(.*)`)
	matches := imageFinder.FindStringSubmatch(string(out))
	if len(matches) < 2 {
		log.Warn().Msgf(L("no UYUNI.*_IMAGE defined in %s systemd service"), service)
		return ""
	}
	return matches[1]
}

// GetImageVirtualSize returns the size of the image with its layers.
func GetImageVirtualSize(name string) (size int64, err error) {
	out, err := utils.NewRunner("podman", "inspect", "--format", "{{.VirtualSize}}", name).
		Log(zerolog.DebugLevel).
		Exec()
	if err != nil {
		return
	}
	sizeStr := strings.TrimSpace(string(out))
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	return
}

// DeleteVolume deletes a podman volume based on its name.
// If dryRun is set to true, nothing will be done, only messages logged to explain what would happen.
func DeleteVolume(name string, dryRun bool) error {
	exists := IsVolumePresent(name)
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
				basePath, errBasePath := GetPodmanVolumeBasePath()
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

// ExportVolume exports a podman volume based on its name to the specified targed directory.
// outputDir option expects already existing directory.
// If dryRun is set to true, only messages will be logged to explain what would happen.
func ExportVolume(name string, outputDir string, dryRun bool) error {
	exists := IsVolumePresent(name)
	if exists {
		outputFile := path.Join(outputDir, name+".tar")
		exportCommand := []string{"podman", "volume", "export", "-o", outputFile, name}
		if dryRun {
			log.Info().Msgf(L("Would run %s"), strings.Join(exportCommand, " "))
			return nil
		}
		log.Info().Msgf(L("Run %s"), strings.Join(exportCommand, " "))
		if err := runCmd(exportCommand[0], exportCommand[1:]...); err != nil {
			return utils.Errorf(err, L("Failed to export volume %s"), name)
		}
		if err := utils.CreateChecksum(outputFile); err != nil {
			return utils.Errorf(err, L("Failed to write checksum of volume %[1]s to the %[2]s"), name, outputFile+".sha256sum")
		}
	}
	return nil
}

// ImportVolume imports a podman volume from provided volumePath.
// If dryRun is set to true, only messages will be logged to exmplain what would happen.
func ImportVolume(name string, volumePath string, skipVerify bool, dryRun bool) error {
	createCommand := []string{"podman", "volume", "create", "--ignore", name}

	basePath, err := GetPodmanVolumeBasePath()
	if err != nil {
		log.Debug().Msg("cannot get base volume path")
		return err
	}
	importCommand := []string{"tar", "xf", volumePath, "-C", path.Join(basePath, name, "_data")}

	if dryRun {
		log.Info().Msgf(L("Would run %s"), strings.Join(importCommand, " "))
		return nil
	}
	if !skipVerify {
		if err := utils.ValidateChecksum(volumePath); err != nil {
			return utils.Errorf(err, L("Checksum does not match for volume %s"), volumePath)
		}
	}
	if err := runCmd(createCommand[0], createCommand[1:]...); err != nil {
		return utils.Errorf(err, L("Failed to precreate empty volume %s"), name)
	}
	log.Info().Msgf(L("Run %s"), strings.Join(importCommand, " "))
	if err := runCmd(importCommand[0], importCommand[1:]...); err != nil {
		return utils.Errorf(err, L("Failed to import volume %s"), name)
	}
	return nil
}

func IsVolumePresent(volume string) bool {
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

// GetPodmanVolumeBasePath returns the path to all volumes on the host system.
func GetPodmanVolumeBasePath() (string, error) {
	cmd := exec.Command("podman", "system", "info", "--format={{ .Store.VolumePath }}")
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// GetVolumeMountPoint returns the path to the volume mount point on the host system.
// This shouldn't be confused with GetPodmanVolumeBasePath() that returns the path to the folder containing all volumes.
func GetVolumeMountPoint(name string) (path string, err error) {
	out, err := utils.NewRunner("podman", "volume", "inspect", "--format", "{{.Mountpoint}}", name).
		Log(zerolog.DebugLevel).
		Exec()
	if err != nil {
		return
	}
	path = strings.TrimSpace(string(out))
	return
}

// Inspect check values on a given image and deploy.
func Inspect(
	serverImage types.ImageFlags,
	pgsqlImage types.ImageFlags,
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

	authFile, cleaner, err := PodmanLogin(hostData, scc, serverImage)
	if err != nil {
		return nil, utils.Errorf(err, L("failed to login to %s"), serverImage.RegistryFQDN)
	}
	defer cleaner()

	podmanArgs := []string{
		"-v", scriptDir + ":" + utils.InspectContainerDirectory,
		"--security-opt", "label=disable",
	}

	preparedImage, err := PrepareImage(authFile, serverImage.Name, serverImage.PullPolicy, true)
	if err != nil {
		return nil, err
	}

	inspector := utils.NewServerInspector(scriptDir)
	if err := inspector.GenerateScript(); err != nil {
		return nil, err
	}
	err = RunContainer("uyuni-inspect", preparedImage.Name, utils.ServerVolumeMounts, podmanArgs,
		[]string{utils.InspectContainerDirectory + "/" + utils.InspectScriptFilename})
	if err != nil {
		return nil, err
	}

	inspectResult, err := inspector.ReadInspectData()
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect data"))
	}

	pgsqlPreparedImage, err := PrepareImage(authFile, pgsqlImage.Name, pgsqlImage.PullPolicy, true)
	if err != nil {
		return nil, err
	}

	dbinspector := utils.NewDBInspector(scriptDir)
	if err := dbinspector.GenerateScript(); err != nil {
		return nil, err
	}

	err = RunContainer("uyuni-db-inspect", pgsqlPreparedImage.Name, utils.PgsqlRequiredVolumeMounts, podmanArgs,
		[]string{utils.InspectContainerDirectory + "/" + utils.InspectScriptFilename})
	if err != nil {
		return nil, err
	}

	dbInspectResult, err := dbinspector.ReadInspectData()
	if err != nil {
		return nil, utils.Errorf(err, L("cannot inspect data"))
	}

	inspectResult.DBInspectData.ImagePgVersion = dbInspectResult.ImagePgVersion

	return inspectResult, err
}
