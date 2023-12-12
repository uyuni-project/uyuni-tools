// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package uninstall

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/podman"
	shared_podman "github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForPodman(dryRun bool, purge bool) {

	// Disable the service
	// Check if there is an uyuni-server service

	if err := utils.RunCmd("systemctl", "list-unit-files", "uyuni-server.service"); err != nil {
		log.Debug().Msg("Systemd has no uyuni-server.service unit")
	} else {
		if dryRun {
			log.Info().Msgf("Would run systemctl disable --now uyuni-server")
			log.Debug().Msgf("Would remove %s", podman.ServicePath)
		} else {
			log.Debug().Msg("Disable uyuni-server service")
			// disable server
			err := utils.RunCmd("systemctl", "disable", "--now", "uyuni-server")
			if err != nil {
				log.Error().Err(err).Msg("Failed to disable server")
			}

			// Remove the service unit
			log.Debug().Msgf("Remove %s", podman.ServicePath)
			if err := os.Remove(podman.ServicePath); err != nil {
				log.Error().Err(err).Msg("Failed to remove uyuni-server.service")
			}
		}
	}

	// Force stop the pod
	containerName := shared_podman.ServerContainerName
	if out, _ := utils.RunCmdOutput(zerolog.DebugLevel, "podman", "ps", "-a", "-q", "-f", "name="+containerName); len(out) > 0 {
		if dryRun {
			log.Debug().Msgf("Would run podman kill %s for container id: %s", containerName, out)
			log.Debug().Msgf("Would run podman remove %s for container id: %s", containerName, out)
		} else {
			log.Debug().Msgf("Run podman kill %s for container id: %s", containerName, out)
			err := utils.RunCmd("podman", "kill", containerName)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to kill the server")

				log.Debug().Msgf("Run podman remove %s for container id: %s", containerName, out)
				err = utils.RunCmd("podman", "rm", containerName)
				if err != nil {
					log.Debug().Err(err).Msg("Error removing container")
				}
			}
		}
	} else {
		log.Debug().Msg("Container already removed")
	}

	// Remove the volumes
	if purge {
		volumes := []string{"cgroup"}
		for volume := range utils.VOLUMES {
			volumes = append(volumes, volume)
		}
		for _, volume := range volumes {
			cmd := exec.Command("podman", "volume", "exists", volume)
			cmd.Run()
			if cmd.ProcessState.ExitCode() == 0 {
				if dryRun {
					log.Debug().Msgf("Would run podman volume rm %s", volume)
				} else {
					errorMessage := fmt.Sprintf("Failed to remove volume %s", volume)
					err := utils.RunCmd("podman", "volume", "rm", volume)
					if err != nil {
						log.Error().Err(err).Msg(errorMessage)
					}
				}
			}
		}
		log.Debug().Msg("All volumes removed")
	}

	// Remove the network
	err := utils.RunCmd("podman", "network", "exists", "uyuni")
	if err != nil {
		log.Info().Msgf("Network uyuni already removed")
	} else {
		if dryRun {
			log.Info().Msgf("Would run podman network rm uyuni")
		} else {
			err := utils.RunCmd("podman", "network", "rm", "uyuni")
			if err != nil {
				log.Error().Msg("Failed to remove network uyuni")
			} else {
				log.Debug().Msg("Network removed")
			}
		}
	}

	// Reload systemd daemon
	if dryRun {
		log.Info().Msg("Would run systemctl reset-failed")
		log.Info().Msg("Would run systemctl daemon-reload")
	} else {
		err := utils.RunCmd("systemctl", "reset-failed")
		if err != nil {
			log.Error().Err(err).Msg("Failed to reset-failed systemd")
		}
		err = utils.RunCmd("systemctl", "daemon-reload")
		if err != nil {
			log.Error().Err(err).Msg("Failed to reload systemd daemon")
		}
	}
}
