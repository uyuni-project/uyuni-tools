// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

// Systemd is an interface providing systemd operations.
type Systemd interface {

	// HasService returns if a systemd service is installed.
	// name is the name of the service without the '.service' part.
	HasService(name string) bool

	// ServiceIsEnabled returns if a service is enabled.
	// name is the name of the service without the '.service' part.
	ServiceIsEnabled(name string) bool

	// EnableService enables and starts a systemd service.
	EnableService(name string) error

	// DisableService disables a service.
	// name is the name of the service without the '.service' part.
	DisableService(name string) error

	// UninstallService stops and remove a systemd service.
	// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
	UninstallService(name string, dryRun bool)

	// ReloadDaemon resets the failed state of services and reload the systemd daemon.
	// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
	ReloadDaemon(dryRun bool) error

	// IsServiceRunning returns whether the systemd service is started or not.
	IsServiceRunning(name string) bool

	// RestartService restarts the systemd service.
	RestartService(name string) error

	// StartService starts the systemd service.
	StartService(name string) error

	// StopService starts the systemd service.
	StopService(name string) error

	// Scales a templated systemd service to the requested number of replicas.
	// name is the name of the service without the '.service' part.
	ScaleService(replicas int, name string) error

	// CurrentReplicaCount returns the current enabled replica count for a template service
	// name is the name of the service without the '.service' part.
	CurrentReplicaCount(name string) int

	// UninstallInstantiatedService stops and remove an instantiated systemd service.
	// If dryRun is set to true, nothing happens but messages are logged to explain what would be done.
	UninstallInstantiatedService(name string, dryRun bool)

	// StartInstantiated starts all replicas.
	StartInstantiated(name string) error

	// RestartInstantiated restarts all replicas.
	RestartInstantiated(name string) error

	// StopInstantiated stops all replicas.
	StopInstantiated(name string) error

	// GetServicesFromSystemdFiles return the uyuni enabled services as string list.
	GetServicesFromSystemdFiles(systemdFileList string) []string
}
