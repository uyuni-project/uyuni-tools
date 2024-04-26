// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const servicesPath = "/etc/systemd/system/"

// Name of the systemd service for the server.
const ServerService = "uyuni-server"

// Name of the systemd service for the coco attestation container.
const ServerAttestationService = "uyuni-server-attestation"

// Name of the systemd service for the proxy.
const ProxyService = "uyuni-proxy-pod"

// HasService returns if a systemd service is installed.
// name is the name of the service without the '.service' part.
func HasService(name string) bool {
	err := utils.RunCmd("systemctl", "list-unit-files", name+".service")
	return err == nil
}

// GetServicePath return the path for a given service.
func GetServicePath(name string) string {
	return path.Join(servicesPath, name+".service")
}

// UninstallService stops and remove a systemd service.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func UninstallService(name string, dryRun bool) {
	servicePath := GetServicePath(name)
	if !HasService(name) {
		log.Info().Msgf(L("Systemd has no %s.service unit"), name)
	} else {
		if dryRun {
			log.Info().Msgf(L("Would run %s"), "systemctl disable --now "+name)
			log.Info().Msgf(L("Would remove %s"), servicePath)
		} else {
			log.Info().Msgf(L("Disable %s service"), name)
			// disable server
			err := utils.RunCmd("systemctl", "disable", "--now", name)
			if err != nil {
				log.Error().Err(err).Msgf(L("Failed to disable %s service"), name)
			}

			// Remove the service unit
			log.Info().Msgf(L("Remove %s"), servicePath)
			if err := os.Remove(servicePath); err != nil {
				log.Error().Err(err).Msgf(L("Failed to remove %s.service file"), name)
			}
		}
	}
}

// ReloadDaemon resets the failed state of services and reload the systemd daemon.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func ReloadDaemon(dryRun bool) error {
	if dryRun {
		log.Info().Msgf(L("Would run %s"), "systemctl reset-failed")
		log.Info().Msgf(L("Would run %s"), "systemctl daemon-reload")
	} else {
		err := utils.RunCmd("systemctl", "reset-failed")
		if err != nil {
			return errors.New(L("failed to reset-failed systemd"))
		}
		err = utils.RunCmd("systemctl", "daemon-reload")
		if err != nil {
			return errors.New(L("failed to reload systemd daemon"))
		}
	}
	return nil
}

// IsServiceRunning returns whether the systemd service is started or not.
func IsServiceRunning(service string) bool {
	cmd := exec.Command("systemctl", "is-active", "-q", service)
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

// RestartService restarts the systemd service.
func RestartService(service string) error {
	if err := utils.RunCmd("systemctl", "restart", service); err != nil {
		return fmt.Errorf(L("failed to restart systemd %s.service: %s"), service, err)
	}
	return nil
}

// StartService starts the systemd service.
func StartService(service string) error {
	if err := utils.RunCmd("systemctl", "start", service); err != nil {
		return fmt.Errorf(L("failed to start systemd %s.service: %s"), service, err)
	}
	return nil
}

// StopService starts the systemd service.
func StopService(service string) error {
	if err := utils.RunCmd("systemctl", "stop", service); err != nil {
		return fmt.Errorf(L("failed to stop systemd %s.service: %s"), service, err)
	}
	return nil
}

// EnableService enables and starts a systemd service.
func EnableService(service string) error {
	if err := utils.RunCmd("systemctl", "enable", "--now", service); err != nil {
		return fmt.Errorf(L("failed to enable %s systemd service: %s"), service, err)
	}
	return nil
}

// Create new systemd service configuration file.
func GenerateSystemdConfFile(serviceName string, section string, body string) error {
	systemdFilePath := GetServicePath(serviceName)

	systemdConfFolder := systemdFilePath + ".d"
	if err := os.MkdirAll(systemdConfFolder, 0750); err != nil {
		return fmt.Errorf(L("failed to create %s folder: %s"), systemdConfFolder, err)
	}
	systemdConfFilePath := path.Join(systemdConfFolder, section+".conf")

	content := []byte("[" + section + "]" + "\n" + body + "\n")
	if err := os.WriteFile(systemdConfFilePath, content, 0644); err != nil {
		return fmt.Errorf(L("cannot write %s file: %s"), systemdConfFilePath, err)
	}

	return nil
}
