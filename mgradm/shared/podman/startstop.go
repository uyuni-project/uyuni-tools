// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func StartServices() error {
	return utils.JoinErrors(
		systemd.StartService(podman.DBService),
		systemd.StartInstantiated(podman.ServerAttestationService),
		systemd.StartInstantiated(podman.EventProcessorService),
		systemd.StartInstantiated(podman.HubXmlrpcService),
		systemd.StartInstantiated(podman.SalineService),
		systemd.StartService(podman.ServerService),
	)
}

func StopServices() error {
	return utils.JoinErrors(
		systemd.StopInstantiated(podman.ServerAttestationService),
		systemd.StopInstantiated(podman.EventProcessorService),
		systemd.StopInstantiated(podman.HubXmlrpcService),
		systemd.StopInstantiated(podman.SalineService),
		systemd.StopService(podman.ServerService),
		systemd.StopService(podman.DBService),
	)
}
