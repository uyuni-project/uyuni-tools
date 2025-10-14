// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

var servicesPath = "/etc/systemd/system/"

// ServerService is the name of the systemd service for the server.
const ServerService = "uyuni-server"

// DBService is the name of the systemd service for the database container.
const DBService = "uyuni-db"

// ServerAttestationService is the name of the systemd service for the coco attestation container.
const ServerAttestationService = "uyuni-server-attestation"

// HubXmlrpcService is the name of the systemd service for the Hub XMLRPC container.
const HubXmlrpcService = "uyuni-hub-xmlrpc"

// SalineService is the name of the systemd service for the saline container.
const SalineService = "uyuni-saline"

// EventProcessorService is the name of the systemd service for the Salt Event Processor container.
const EventProcessorService = "uyuni-event-processor"

// ProxyService is the name of the systemd service for the proxy.
const ProxyService = "uyuni-proxy-pod"

// CustomConf is the name of the custom configuration file of services.
const CustomConf = "custom.conf"

// Interface to perform systemd calls.
// This is not meant to be used elsewhere than in the SystemdImpl class and the unit tests.
type SystemdDriver interface {
	// HasService returns if a systemd service is installed.
	// name is the name of the service without the '.service' part.
	HasService(name string) bool

	// ServiceIsEnabled returns if a service is enabled
	// name is the name of the service without the '.service' part.
	ServiceIsEnabled(name string) bool

	// DisableService disables a service
	// name is the name of the service without the '.service' part.
	DisableService(name string) error

	// ReloadDaemon resets the failed state of services and reload the systemd daemon.
	// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
	ReloadDaemon() error

	// IsServiceRunning returns whether the systemd service is started or not.
	IsServiceRunning(service string) bool

	// RestartService restarts the systemd service.
	RestartService(service string) error

	// StartService starts the systemd service.
	StartService(service string) error

	// StopService starts the systemd service.
	StopService(service string) error

	// EnableService enables and starts a systemd service.
	EnableService(service string) error

	// GetServiceProperty returns the value of a systemd service property.
	GetServiceProperty(service string, property string) (string, error)
}

type systemdDriverImpl struct {
}

// HasService returns if a systemd service is installed.
// name is the name of the service without the '.service' part.
func (d *systemdDriverImpl) HasService(name string) bool {
	err := utils.RunCmd("systemctl", "list-unit-files", name+".service")
	return err == nil
}

// ServiceIsEnabled returns if a service is enabled
// name is the name of the service without the '.service' part.
func (d *systemdDriverImpl) ServiceIsEnabled(name string) bool {
	err := utils.RunCmd("systemctl", "is-enabled", name+".service")
	return err == nil
}

// DisableService disables a service
// name is the name of the service without the '.service' part.
func (d *systemdDriverImpl) DisableService(name string) error {
	if !d.ServiceIsEnabled(name) {
		log.Debug().Msgf("%s is already disabled.", name)
		return nil
	}
	if err := utils.RunCmd("systemctl", "disable", "--now", name); err != nil {
		return utils.Errorf(err, L("failed to disable %s systemd service"), name)
	}
	return nil
}

// ReloadDaemon resets the failed state of services and reload the systemd daemon.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func (d *systemdDriverImpl) ReloadDaemon() error {
	err := utils.RunCmd("systemctl", "reset-failed")
	if err != nil {
		return errors.New(L("failed to reset-failed systemd"))
	}
	err = utils.RunCmd("systemctl", "daemon-reload")
	if err != nil {
		return errors.New(L("failed to reload systemd daemon"))
	}
	return nil
}

// IsServiceRunning returns whether the systemd service is started or not.
func (d *systemdDriverImpl) IsServiceRunning(service string) bool {
	cmd := exec.Command("systemctl", "is-active", "-q", service)
	if err := cmd.Run(); err != nil {
		return false
	}
	return cmd.ProcessState.ExitCode() == 0
}

// RestartService restarts the systemd service.
func (d *systemdDriverImpl) RestartService(service string) error {
	if err := utils.RunCmd("systemctl", "restart", service); err != nil {
		return utils.Errorf(err, L("failed to restart systemd %s.service"), service)
	}
	return nil
}

// StartService starts the systemd service.
func (d *systemdDriverImpl) StartService(service string) error {
	if err := utils.RunCmd("systemctl", "start", service); err != nil {
		return utils.Errorf(err, L("failed to start systemd %s.service"), service)
	}
	return nil
}

// StopService starts the systemd service.
func (d *systemdDriverImpl) StopService(service string) error {
	if err := utils.RunCmd("systemctl", "stop", service); err != nil {
		return utils.Errorf(err, L("failed to stop systemd %s.service"), service)
	}
	return nil
}

// EnableService enables and starts a systemd service.
func (d *systemdDriverImpl) EnableService(service string) error {
	if d.ServiceIsEnabled(service) {
		log.Debug().Msgf("%s is already enabled.", service)
		return nil
	}
	if err := utils.RunCmd("systemctl", "enable", "--now", service); err != nil {
		return utils.Errorf(err, L("failed to enable %s systemd service"), service)
	}
	return nil
}

func (d *systemdDriverImpl) GetServiceProperty(service string, property string) (string, error) {
	serviceName := service
	if strings.HasSuffix(service, "@") {
		serviceName = service + "0"
	}
	out, err := newRunner("systemctl", "show", "-p", property, serviceName).Exec()
	if err != nil {
		return "", utils.Errorf(err, L("Failed to get the %[1]s property from %[2]s service"), property, service)
	}
	return strings.TrimPrefix(strings.TrimSpace(string(out)), property+"="), nil
}

// NewSystemd returns a new Systemd instance.
func NewSystemd() Systemd {
	driver := systemdDriverImpl{}
	return SystemdImpl{
		driver: &driver,
	}
}

// NewSystemdWithDriver returns a new Systemd instance with a custom driver.
func NewSystemdWithDriver(driver SystemdDriver) Systemd {
	return SystemdImpl{
		driver: driver,
	}
}

// SystemdImpl implements the Systemd interface.
type SystemdImpl struct {
	// driver actually calls the systemd executable.
	driver SystemdDriver
}

// HasService returns if a systemd service is installed.
// name is the name of the service without the '.service' part.
func (s SystemdImpl) HasService(name string) bool {
	return s.driver.HasService(name)
}

// ServiceIsEnabled returns if a service is enabled
// name is the name of the service without the '.service' part.
func (s SystemdImpl) ServiceIsEnabled(name string) bool {
	return s.driver.ServiceIsEnabled(name)
}

// DisableService disables a service
// name is the name of the service without the '.service' part.
func (s SystemdImpl) DisableService(name string) error {
	return s.driver.DisableService(name)
}

// GetServicePath return the path for a given service.
func GetServicePath(name string) string {
	return path.Join(servicesPath, name+".service")
}

func (s SystemdImpl) GetServiceProperty(service string, property string) (string, error) {
	return s.driver.GetServiceProperty(service, property)
}

// GetServiceConfFolder return the conf folder for systemd services.
func GetServiceConfFolder(name string) string {
	return path.Join(servicesPath, name+".service.d")
}

// GetServiceConfPath return the path for generated.conf file.
func GetServiceConfPath(name string) string {
	return path.Join(GetServiceConfFolder(name), "generated.conf")
}

// UninstallService stops and remove a systemd service.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func (s SystemdImpl) UninstallService(name string, dryRun bool) {
	if !s.HasService(name) {
		log.Info().Msgf(L("Systemd has no %s.service unit"), name)
	} else {
		if dryRun {
			log.Info().Msgf(L("Would run %s"), "systemctl disable --now "+name)
		} else {
			log.Info().Msgf(L("Disable %s service"), name)
			err := s.DisableService(name)
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
		uninstallServiceFiles(name, dryRun)
	}
}

func uninstallServiceFiles(name string, dryRun bool) {
	servicePath := GetServicePath(name)
	serviceConfFolder := GetServiceConfFolder(name)

	if dryRun {
		log.Info().Msgf(L("Would remove %s"), servicePath)
	} else {
		// Remove the service unit
		log.Info().Msgf(L("Remove %s"), servicePath)
		if err := os.Remove(servicePath); err != nil {
			log.Error().Err(err).Msgf(L("Failed to remove %s.service file"), name)
		}
	}

	if utils.FileExists(serviceConfFolder) {
		confPaths := []string{
			GetServiceConfPath(name),
			path.Join(serviceConfFolder, "Service.conf"),
		}
		for _, confPath := range confPaths {
			if utils.FileExists(confPath) {
				if dryRun {
					log.Info().Msgf(L("Would remove %s"), confPath)
				} else {
					log.Info().Msgf(L("Remove %s"), confPath)
					if err := os.Remove(confPath); err != nil {
						log.Error().Err(err).Msgf(L("Failed to remove %s file"), confPath)
					}
				}
			}
		}

		if dryRun {
			log.Info().Msgf(L("Would remove %s if empty"), serviceConfFolder)
		} else {
			if utils.IsEmptyDirectory(serviceConfFolder) {
				log.Debug().Msgf("Removing %s folder, since it's empty", serviceConfFolder)
				_ = utils.RemoveDirectory(serviceConfFolder)
			} else {
				log.Warn().Msgf(
					L("%s folder contains file created by the user. Please remove them when uninstallation is completed."),
					serviceConfFolder,
				)
			}
		}
	}
}

// UninstallInstantiatedService stops and remove an instantiated systemd service.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func (s SystemdImpl) UninstallInstantiatedService(name string, dryRun bool) {
	if dryRun {
		log.Info().Msgf(L("Would scale %s to 0 replicas"), name)
	} else {
		if err := s.ScaleService(0, name); err != nil {
			log.Error().Err(err).Msgf(L("Failed to disable %s service"), name)
		}
	}

	uninstallServiceFiles(name+"@", dryRun)
}

// ReloadDaemon resets the failed state of services and reload the systemd daemon.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func (s SystemdImpl) ReloadDaemon(dryRun bool) error {
	if dryRun {
		log.Info().Msgf(L("Would run %s"), "systemctl reset-failed")
		log.Info().Msgf(L("Would run %s"), "systemctl daemon-reload")
	} else {
		return s.driver.ReloadDaemon()
	}
	return nil
}

// IsServiceRunning returns whether the systemd service is started or not.
func (s SystemdImpl) IsServiceRunning(service string) bool {
	return s.driver.IsServiceRunning(service)
}

// RestartService restarts the systemd service.
func (s SystemdImpl) RestartService(service string) error {
	return s.driver.RestartService(service)
}

// StartService starts the systemd service.
func (s SystemdImpl) StartService(service string) error {
	return s.driver.StartService(service)
}

// StopService starts the systemd service.
func (s SystemdImpl) StopService(service string) error {
	return s.driver.StopService(service)
}

// EnableService enables and starts a systemd service.
func (s SystemdImpl) EnableService(service string) error {
	return s.driver.EnableService(service)
}

// StartInstantiated starts all replicas.
func (s SystemdImpl) StartInstantiated(service string) error {
	var errList []error
	for i := 0; i < s.CurrentReplicaCount(service); i++ {
		err := s.StartService(fmt.Sprintf("%s@%d", service, i))
		errList = append(errList, err)
	}
	return utils.JoinErrors(errList...)
}

// RestartInstantiated restarts all replicas.
func (s SystemdImpl) RestartInstantiated(service string) error {
	var errList []error
	for i := 0; i < s.CurrentReplicaCount(service); i++ {
		err := s.RestartService(fmt.Sprintf("%s@%d", service, i))
		errList = append(errList, err)
	}
	return utils.JoinErrors(errList...)
}

// StopInstantiated stops all replicas.
func (s SystemdImpl) StopInstantiated(service string) error {
	var errList []error
	for i := 0; i < s.CurrentReplicaCount(service); i++ {
		err := s.StopService(fmt.Sprintf("%s@%d", service, i))
		errList = append(errList, err)
	}
	return utils.JoinErrors(errList...)
}

// CurrentReplicaCount returns the current enabled replica count for a template service
// name is the name of the service without the '.service' part.
func (s SystemdImpl) CurrentReplicaCount(name string) int {
	count := 0
	for s.ServiceIsEnabled(fmt.Sprintf("%s@%d", name, count)) {
		count++
	}
	return count
}

// ScaleService scales a templated systemd service to the requested number of replicas.
// name is the name of the service without the '.service' part.
func (s SystemdImpl) ScaleService(replicas int, name string) error {
	currentReplicas := s.CurrentReplicaCount(name)
	if currentReplicas == replicas {
		log.Info().Msgf(L("Service %[1]s already has %[2]d replicas."), name, currentReplicas)
		return nil
	}
	log.Info().Msgf(L("Scale %[1]s from %[2]d to %[3]d replicas."), name, currentReplicas, replicas)
	for i := currentReplicas; i < replicas; i++ {
		serviceName := fmt.Sprintf("%s@%d", name, i)
		if err := s.EnableService(serviceName); err != nil {
			return utils.Errorf(err, L("cannot enable service"))
		}
	}
	for i := replicas; i < currentReplicas; i++ {
		serviceName := fmt.Sprintf("%s@%d", name, i)
		if err := s.DisableService(serviceName); err != nil {
			return utils.Errorf(err, L("cannot disable service"))
		}
	}
	return s.RestartInstantiated(name)
}

// Show calls the systemctl show command and returns the output.
func (s SystemdImpl) Show(service string, property string) (string, error) {
	out, err := newRunner("systemctl", "show", "--property", property, service).Log(zerolog.DebugLevel).Exec()
	if err != nil {
		return "", utils.Errorf(err, L("failed to show %[1]s property of %[2] systemd service"), property, service)
	}
	return strings.TrimSpace(string(out)), nil
}

// confHeader is the header for the generated systemd configuration files.
const confHeader = `# This file is generated by mgradm and will be overwritten during upgrades.
# Custom configuration should go in another .conf file in the same folder.

`

// GenerateSystemdConfFile creates a new systemd service configuration file (e.g. Service.conf).
func GenerateSystemdConfFile(serviceName string, filename string, body string, withHeader bool) error {
	systemdFilePath := GetServicePath(serviceName)

	systemdConfFolder := systemdFilePath + ".d"
	if err := os.MkdirAll(systemdConfFolder, 0750); err != nil {
		return utils.Errorf(err, L("failed to create %s folder"), systemdConfFolder)
	}
	systemdConfFilePath := path.Join(systemdConfFolder, filename)

	header := ""
	if withHeader {
		header = confHeader
	}
	content := []byte(fmt.Sprintf("%s[Service]\n%s\n", header, body))
	if err := os.WriteFile(systemdConfFilePath, content, 0640); err != nil {
		return utils.Errorf(err, L("cannot write %s file"), systemdConfFilePath)
	}

	return nil
}

// CleanSystemdConfFile separates the Service.conf file once generated into generated.conf and custom.conf.
func CleanSystemdConfFile(serviceName string) error {
	systemdFilePath := GetServicePath(serviceName) + ".d"
	oldConfPath := path.Join(systemdFilePath, "Service.conf")

	// The first containerized release generated a Service.conf where the image and the configuration
	// where stored. This had the side effect to remove the conf at upgrade time.
	// If this file exists split it in two:
	// - generated.conf with the image
	// - custom.conf with everything that shouldn't be touched at upgrade
	if utils.FileExists(oldConfPath) {
		content := string(utils.ReadFile(oldConfPath))
		lines := strings.Split(content, "\n")

		generated := ""
		custom := ""
		hasCustom := false

		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "Environment=UYUNI_IMAGE=") {
				generated = generated + trimmedLine
			} else {
				custom = custom + trimmedLine + "\n"
				if trimmedLine != "" && trimmedLine != "[Service]" {
					hasCustom = true
				}
			}
		}

		if generated != "" {
			if err := GenerateSystemdConfFile(serviceName, "generated.conf", generated, true); err != nil {
				return err
			}
		}

		if hasCustom {
			customPath := path.Join(systemdFilePath, CustomConf)
			if err := os.WriteFile(customPath, []byte(custom), 0644); err != nil {
				return utils.Errorf(err, L("failed to write %s file"), customPath)
			}
		}

		if err := os.Remove(oldConfPath); err != nil {
			return utils.Errorf(err, L("failed to remove old %s systemd service configuration file"), oldConfPath)
		}
	}

	return nil
}
