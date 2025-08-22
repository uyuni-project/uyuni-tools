// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"errors"
	"fmt"
)

// FakeSystemdDriver is a dummy implementation of the systemd driver for unit tests.
type FakeSystemdDriver struct {
	// Installed is the slice of installed services.
	// Instantiated services are to be listed with the trailing @.
	Installed []string

	// Enabled is the slice of enabled services.
	// All the instances of the instantiated services need to be listed here.
	Enabled []string

	// Running is the slice of running services.
	// All the instances of the instantiated services need to be listed here.
	Running []string

	// DisableServiceErrors maps an error with a service name to mock errors in DisableService.
	DisableServiceErrors map[string]error

	// EnableServiceErrors maps an error with a service name to mock errors in EnableService.
	EnableServiceErrors map[string]error

	// ReloadDaemonError is the error to return in ReloadDaemon.
	ReloadDaemonError error

	// RestartServiceErrors maps an error with a service name to mock errors in RestartService.
	RestartServiceErrors map[string]error

	// StartServiceErrors maps an error with a service name to mock errors in StartService.
	StartServiceErrors map[string]error

	// StopServiceErrors maps an error with a service name to mock errors in StopService.
	StopServiceErrors map[string]error

	// ServiceProperties maps all the properties of each service.
	ServiceProperties map[string]map[string]string
}

// HasService returns if a systemd service is installed.
// name is the name of the service without the '.service' part.
func (d *FakeSystemdDriver) HasService(name string) bool {
	return contains(d.Installed, name)
}

// ServiceIsEnabled returns if a service is enabled
// name is the name of the service without the '.service' part.
func (d *FakeSystemdDriver) ServiceIsEnabled(name string) bool {
	return contains(d.Enabled, name)
}

// DisableService disables a service
// name is the name of the service without the '.service' part.
func (d *FakeSystemdDriver) DisableService(name string) error {
	if !d.ServiceIsEnabled(name) {
		return fmt.Errorf("%s service is not enabled", name)
	}
	err := d.DisableServiceErrors[name]
	if err == nil {
		d.Enabled = deleteItems(d.Enabled, name)
	}
	return err
}

// EnableService enables and starts a systemd service.
func (d *FakeSystemdDriver) EnableService(service string) error {
	err := d.EnableServiceErrors[service]
	if err == nil && !contains(d.Enabled, service) {
		d.Enabled = append(d.Enabled, service)
	}
	return err
}

// ReloadDaemon resets the failed state of services and reload the systemd daemon.
// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
func (d *FakeSystemdDriver) ReloadDaemon() error {
	return d.ReloadDaemonError
}

// IsServiceRunning returns whether the systemd service is started or not.
func (d *FakeSystemdDriver) IsServiceRunning(service string) bool {
	return contains(d.Running, service)
}

// RestartService restarts the systemd service.
func (d *FakeSystemdDriver) RestartService(service string) error {
	if !d.ServiceIsEnabled(service) {
		return fmt.Errorf("%s service is not enabled", service)
	}
	err := d.RestartServiceErrors[service]
	// Same implementation than start, may be this needs to be enhanced for unit tests observability
	if err == nil && !contains(d.Running, service) {
		d.Running = append(d.Running, service)
	}
	return err
}

// StartService starts the systemd service.
func (d *FakeSystemdDriver) StartService(service string) error {
	if !d.ServiceIsEnabled(service) {
		return fmt.Errorf("%s service is not enabled", service)
	}
	err := d.StartServiceErrors[service]
	if err == nil && !contains(d.Running, service) {
		d.Running = append(d.Running, service)
	}
	return err
}

// StopService starts the systemd service.
func (d *FakeSystemdDriver) StopService(service string) error {
	if !d.ServiceIsEnabled(service) {
		return fmt.Errorf("%s service is not enabled", service)
	}
	err := d.StopServiceErrors[service]
	if err == nil {
		d.Running = deleteItems(d.Running, service)
	}
	return err
}

// GetServiceProperty gets the value from the ServiceProperties structure.
// An error is returned if either the service or property doesn't exist.
func (d *FakeSystemdDriver) GetServiceProperty(service string, property string) (string, error) {
	properties, exists := d.ServiceProperties[service]
	if !exists {
		return "", errors.New("no such service")
	}
	value, exists := properties[property]
	if !exists {
		return "", errors.New("no such property")
	}
	return value, nil
}

// deleteItems removes all items equal to needle in the slice.
func deleteItems(slice []string, needle string) []string {
	cleaned := []string{}
	for _, item := range slice {
		if item != needle {
			cleaned = append(cleaned, item)
		}
	}
	return cleaned
}
