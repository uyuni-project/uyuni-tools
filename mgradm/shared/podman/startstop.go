// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"errors"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func StartServices() error {
	return errors.Join(
		systemd.StartService(podman.DBService),
		systemd.StartInstantiated(podman.ServerAttestationService),
		systemd.StartInstantiated(podman.HubXmlrpcService),
		systemd.StartService(podman.ServerService),
	)
}

func StopServices() error {
	return errors.Join(
		systemd.StopInstantiated(podman.ServerAttestationService),
		systemd.StopInstantiated(podman.HubXmlrpcService),
		systemd.StopService(podman.ServerService),
		systemd.StopService(podman.DBService),
	)
}
