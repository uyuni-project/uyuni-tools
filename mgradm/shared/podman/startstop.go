// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
)

func StartServices() error {
	WarnIfServicesDisabled(systemd)

	var dbErr error
	if systemd.HasService(podman.DBService) {
		dbErr = systemd.StartService(podman.DBService)
	}
	errs := errors.Join(
		dbErr,
		systemd.StartInstantiated(podman.ServerAttestationService),
		systemd.StartInstantiated(podman.HubXmlrpcService),
		systemd.StartInstantiated(podman.SalineService),
		systemd.StartService(podman.ServerService),
	)

	if systemd.ServiceIsEnabled(podman.TFTPService) {
		errs = errors.Join(errs, systemd.StartService(podman.TFTPService))
	}

	if systemd.HasService(podman.SalineService + "@") {
		errs = errors.Join(errs, systemd.StartInstantiated(podman.SalineService))
	}

	return errs
}

func StopServices() error {
	WarnIfServicesDisabled(systemd)

	errs := errors.Join(
		systemd.StopInstantiated(podman.ServerAttestationService),
		systemd.StopInstantiated(podman.HubXmlrpcService),
		systemd.StopInstantiated(podman.SalineService),
		systemd.StopService(podman.ServerService),
	)

	if systemd.HasService(podman.DBService) {
		errs = errors.Join(errs, systemd.StopService(podman.DBService))
	}

	if systemd.HasService(podman.SalineService + "@") {
		errs = errors.Join(errs, systemd.StopInstantiated(podman.SalineService))
	}

	if systemd.ServiceIsEnabled(podman.TFTPService) {
		errs = errors.Join(errs, systemd.StopService(podman.TFTPService))
	}
	return errs
}
