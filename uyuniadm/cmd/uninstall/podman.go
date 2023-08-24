package uninstall

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	"github.com/uyuni-project/uyuni-tools/uyuniadm/shared/podman"
)

func uninstallForPodman(globalFlags *types.GlobalFlags, dryRun bool, purge bool) {
	// Check if there is an uyuni-server service
	if err := exec.Command("systemctl", "list-unit-files", "uyuni-server.service").Run(); err != nil {
		log.Fatalln("Systemd has no uyuni-server.service unit, nothing to uninstall")
	}

	// Force stop the pod
	if out, _ := exec.Command("podman", "ps", "-q", "-f", "name=uyuni-server").Output(); len(out) > 0 {
		if dryRun {
			log.Println("Would run podman kill uyuni-server")
		} else {
			utils.RunCmd("podman", []string{"kill", "uyuni-server"}, "Failed to kill the server", globalFlags.Verbose)
		}
	}

	// Disable the service
	if dryRun {
		log.Println("Would run systemctl disable --now uyuni-server")
	} else {
		utils.RunCmd("systemctl", []string{"disable", "--now", "uyuni-server"}, "Failed to disable server", globalFlags.Verbose)
	}

	// Remove the volumes
	if purge {
		for volume := range utils.VOLUMES {
			if dryRun {
				log.Printf("Would run podman volume rm %s\n", volume)
			} else {
				errorMessage := fmt.Sprintf("Failed to remove volume %s", volume)
				utils.RunCmd("podman", []string{"volume", "rm", volume}, errorMessage, globalFlags.Verbose)
			}
		}

		if dryRun {
			log.Println("Would run podman volume rm cgroup")
		} else {
			utils.RunCmd("podman", []string{"volume", "rm", "cgroup"}, "Failed to remove volume cgroup", globalFlags.Verbose)
		}
	}

	// Remove the service unit
	if dryRun {
		log.Printf("Woud remove %s\n", podman.ServicePath)
	} else {
		if globalFlags.Verbose {
			log.Printf("Remove %s\n", podman.ServicePath)
		}
		os.Remove(podman.ServicePath)
	}

	// Remove the network
	cmd := exec.Command("podman", "network", "exists", "uyuni")
	err := cmd.Run()
	if err != nil {
		log.Println("Network uyuni already removed")
		return
	}
	if dryRun {
		log.Println("Would run podman network rm uyuni")
	} else {
		utils.RunCmd("podman", []string{"network", "rm", "uyuni"}, "Failed to remove network uyuni", globalFlags.Verbose)
	}

	// Reload systemd daemon
	if dryRun {
		log.Println("Would run systemctl daemon-reload")
	} else {
		utils.RunCmd("systemctl", []string{"daemon-reload"}, "Failed to reload systemd daemon", globalFlags.Verbose)
	}
}
