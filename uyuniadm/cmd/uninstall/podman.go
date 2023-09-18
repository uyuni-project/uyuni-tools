package uninstall

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
)

func uninstallForPodman(globalFlags *types.GlobalFlags, dryRun bool, purge bool) {

	// Disable the service
	// Check if there is an uyuni-server service

	if err := utils.RunRawCmd("systemctl", []string{"list-unit-files", "uyuni-server.service"}, true); err != nil {
		log.Debug().Msg("Systemd has no uyuni-server.service unit")
	} else {
		if dryRun {
			log.Info().Msgf("Would run systemctl disable --now uyuni-server")
			log.Debug().Msgf("Woud remove %s", podman.ServicePath)
		} else {
			log.Debug().Msg("Desable uyuni-server service")
			// disable server
			err := utils.RunRawCmd("systemctl", []string{"disable", "--now", "uyuni-server"}, true)
			if err != nil {
				log.Error().Err(err).Msg("Failed to disable server")
			}

			// Remove the service unit
			log.Debug().Msgf("Remove %s", podman.ServicePath)
			os.Remove(podman.ServicePath)
		}
	}

	// Force stop the pod
	if out, _ := exec.Command("podman", "ps", "-a", "-q", "-f", "name=uyuni-server").Output(); len(out) > 0 {
		if dryRun {
			log.Debug().Msgf("Would run podman kill uyuni-server for container id: %s", out)
			log.Debug().Msgf("Would run podman remove uyuni-server for container id: %s", out)
		} else {
			log.Debug().Msgf("Run podman kill uyuni-server for container id: %s", out)
			err := utils.RunRawCmd("podman", []string{"kill", "uyuni-server"}, true)
			if err != nil {
				log.Debug().Err(err).Msg("Failed to kill the server")

				log.Debug().Msgf("Run podman remove uyuni-server for container id: %s", out)
				err = utils.RunRawCmd("podman", []string{"rm", "uyuni-server"}, true)
				if err != nil {
					log.Debug().Err(err).Msg("Error removing container")
				}
			}
		}
	} else {
		log.Debug().Msg("Container already remove")
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
					err := utils.RunRawCmd("podman", []string{"volume", "rm", volume}, true)
					if err != nil {
						log.Error().Err(err).Msg(errorMessage)
					}
				}
			}
		}
	}

	// Remove the network
	err := utils.RunRawCmd("podman", []string{"network", "exists", "uyuni"}, false)
	if err != nil {
		log.Info().Msgf("Network uyuni already removed")
	} else {
		if dryRun {
			log.Info().Msgf("Would run podman network rm uyuni")
		} else {
			err := utils.RunRawCmd("podman", []string{"network", "rm", "uyuni"}, true)
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
		err := utils.RunRawCmd("systemctl", []string{"reset-failed"}, true)
		if err != nil {
			log.Error().Err(err).Msg("Failed to reset-failed systemd")
		}
		err = utils.RunRawCmd("systemctl", []string{"daemon-reload"}, true)
		if err != nil {
			log.Error().Err(err).Msg("Failed to reload systemd daemon")
		}
	}
}
