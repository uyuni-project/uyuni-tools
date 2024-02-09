// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
package types

type VolumeMount struct {
	MountPath string `json:"mountPath,omitempty"`
	Name      string `json:"name,omitempty"`
}

type Container struct {
	Name         string        `json:"name,omitempty"`
	Image        string        `json:"image,omitempty"`
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
}

type PersistentVolumeClaim struct {
	ClaimName string `json:"claimName,omitempty"`
}

type HostPath struct {
	Path string `json:"path,omitempty"`
	Type string `json:"type,omitempty"`
}

type SecretItem struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
}

type Secret struct {
	SecretName string       `json:"secretName,omitempty"`
	Items      []SecretItem `json:"items,omitempty"`
}

type Volume struct {
	Name                  string                 `json:"name,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`
	HostPath              *HostPath              `json:"hostPath,omitempty"`
	Secret                *Secret                `json:"secret,omitempty"`
}

type Spec struct {
	NodeName      string      `json:"nodeName,omitempty"`
	RestartPolicy string      `json:"restartPolicy,omitempty"`
	Containers    []Container `json:"containers,omitempty"`
	Volumes       []Volume    `json:"volumes,omitempty"`
}
type Deployment struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Spec       *Spec  `json:"spec,omitempty"`
}
