// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const servicesPath = "/etc/systemd/system/"

// Name of the systemd service for the server.
const ServerService = "uyuni-server"

// Name of the systemd service for the proxy.
const ProxyService = "uyuni-proxy-pod"

// HasService returns if a systemd service is installed.
// name is the name of the service without the '.service' part.
func HasService(name string) bool {
	err := utils.RunCmd("systemctl", "list-unit-files", name+".service")
	return err != nil
}

func GetServicePath(name string) string {
	return path.Join(servicesPath, name+".service")
}

// UninstallService stops and remove a systemd service.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func UninstallService(name string, dryRun bool) {
	servicePath := GetServicePath(name)
	if HasService(name) {
		log.Info().Msgf("Systemd has no %s.service unit", name)
	} else {
		if dryRun {
			log.Info().Msgf("Would run systemctl disable --now %s", name)
			log.Info().Msgf("Would remove %s", servicePath)
		} else {
			log.Info().Msgf("Disable %s service", name)
			// disable server
			err := utils.RunCmd("systemctl", "disable", "--now", name)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to disable %s service", name)
			}

			// Remove the service unit
			log.Info().Msgf("Remove %s", servicePath)
			if err := os.Remove(servicePath); err != nil {
				log.Error().Err(err).Msgf("Failed to remove %s.service file", name)
			}
		}
	}
}

// ReloadDaemon resets the failed state of services and reload the systemd daemon.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func ReloadDaemon(dryRun bool) {
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

// IsServiceRunning returns whether the systemd service is started or not.
func IsServiceRunning(service string) bool {
	cmd := exec.Command("systemctl", "is-active", "-q", service)
	cmd.Run()
	return cmd.ProcessState.ExitCode() == 0
}

// RestartService restarts the systemd service.
func RestartService(service string) error {
	if err := utils.RunCmd("systemctl", "restart", service); err != nil {
		return fmt.Errorf("failed to restart systemd %s.service: %s", service, err)
	}
	return nil
}

// StartService starts the systemd service.
func StartService(service string) error {
	if err := utils.RunCmd("systemctl", "start", service); err != nil {
		return fmt.Errorf("failed to start systemd %s.service: %s", service, err)
	}
	return nil
}

// StopService starts the systemd service.
func StopService(service string) error {
	if err := utils.RunCmd("systemctl", "stop", service); err != nil {
		return fmt.Errorf("failed to stop systemd %s.service: %s", service, err)
	}
	return nil
}

// EnableService enables and starts a systemd service.
func EnableService(service string) error {
	if err := utils.RunCmd("systemctl", "enable", "--now", service); err != nil {
		return fmt.Errorf("failed to enable %s systemd service: %s", service, err)
	}
	return nil
}
