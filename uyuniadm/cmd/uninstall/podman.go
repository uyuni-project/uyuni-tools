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
	// Check if there is an uyuni-server service
	if err := exec.Command("systemctl", "list-unit-files", "uyuni-server.service").Run(); err != nil {
		log.Fatal().Msg("Systemd has no uyuni-server.service unit, nothing to uninstall")
	}

	// Force stop the pod
	if out, _ := exec.Command("podman", "ps", "-q", "-f", "name=uyuni-server").Output(); len(out) > 0 {
		if dryRun {
			log.Info().Msgf("Would run podman kill uyuni-server")
		} else {
			utils.RunCmd("podman", []string{"kill", "uyuni-server"}, "Failed to kill the server", globalFlags.Verbose)
		}
	}

	// Disable the service
	if dryRun {
		log.Info().Msgf("Would run systemctl disable --now uyuni-server")
	} else {
		utils.RunCmd("systemctl", []string{"disable", "--now", "uyuni-server"}, "Failed to disable server", globalFlags.Verbose)
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
					log.Info().Msgf("Would run podman volume rm %s", volume)
				} else {
					errorMessage := fmt.Sprintf("Failed to remove volume %s", volume)
					utils.RunCmd("podman", []string{"volume", "rm", volume}, errorMessage, globalFlags.Verbose)
				}
			}
		}
	}

	// Remove the service unit
	if dryRun {
		log.Info().Msgf("Woud remove %s", podman.ServicePath)
	} else {
		if globalFlags.Verbose {
			log.Info().Msgf("Remove %s", podman.ServicePath)
		}
		os.Remove(podman.ServicePath)
	}

	// Remove the network
	cmd := exec.Command("podman", "network", "exists", "uyuni")
	err := cmd.Run()
	if err != nil {
		log.Info().Msgf("Network uyuni already removed")
		return
	}
	if dryRun {
		log.Info().Msgf("Would run podman network rm uyuni")
	} else {
		utils.RunCmd("podman", []string{"network", "rm", "uyuni"}, "Failed to remove network uyuni", globalFlags.Verbose)
	}

	// Reload systemd daemon
	if dryRun {
		log.Info().Msg("Would run systemctl daemon-reload")
	} else {
		utils.RunCmd("systemctl", []string{"reset-failed"}, "Failed to reload systemd daemon", globalFlags.Verbose)
		utils.RunCmd("systemctl", []string{"daemon-reload"}, "Failed to reload systemd daemon", globalFlags.Verbose)
	}
}
