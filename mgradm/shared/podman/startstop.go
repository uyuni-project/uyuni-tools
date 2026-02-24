// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func StartServices() error {
	var dbErr error
	if systemd.HasService(podman.DBService) {
		dbErr = systemd.StartService(podman.DBService)
	}
	errs := utils.JoinErrors(
		dbErr,
		systemd.StartInstantiated(podman.ServerAttestationService),
		systemd.StartInstantiated(podman.HubXmlrpcService),
		systemd.StartInstantiated(podman.SalineService),
		systemd.StartService(podman.ServerService),
	)

	if systemd.ServiceIsEnabled(podman.TFTPService) {
		errs = utils.JoinErrors(errs, systemd.StartService(podman.TFTPService))
	}

	if systemd.HasService(podman.SalineService + "@") {
		errs = utils.JoinErrors(errs, systemd.StartInstantiated(podman.SalineService))
	}

	return errs
}

func StopServices() error {
	errs := utils.JoinErrors(
		systemd.StopInstantiated(podman.ServerAttestationService),
		systemd.StopInstantiated(podman.HubXmlrpcService),
		systemd.StopInstantiated(podman.SalineService),
		systemd.StopService(podman.ServerService),
	)

	if systemd.HasService(podman.DBService) {
		errs = utils.JoinErrors(errs, systemd.StopService(podman.DBService))
	}

	if systemd.HasService(podman.SalineService + "@") {
		errs = utils.JoinErrors(errs, systemd.StopInstantiated(podman.SalineService))
	}

	if systemd.ServiceIsEnabled(podman.TFTPService) {
		errs = utils.JoinErrors(errs, systemd.StopService(podman.TFTPService))
	}
	return errs
}
