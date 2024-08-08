// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"github.com/uyuni-project/uyuni-tools/shared/types"
	core "k8s.io/api/core/v1"
)

// ConvertVolumeMounts converts the internal volume mounts into Kubernetes' ones.
func ConvertVolumeMounts(mounts []types.VolumeMount) []core.VolumeMount {
	res := []core.VolumeMount{}

	for _, mount := range mounts {
		converted := core.VolumeMount{
			Name:      mount.Name,
			MountPath: mount.MountPath,
		}
		res = append(res, converted)
	}

	return res
}

// ConvertPortMaps converts the internal port maps to Kubernetes ContainerPorts.
func ConvertPortMaps(ports []types.PortMap) []core.ContainerPort {
	res := []core.ContainerPort{}

	for _, port := range ports {
		protocol := core.ProtocolTCP
		if port.Protocol == "udp" {
			protocol = core.ProtocolUDP
		}
		converted := core.ContainerPort{
			ContainerPort: int32(port.Exposed),
			Protocol:      protocol,
		}
		res = append(res, converted)
	}
	return res
}
