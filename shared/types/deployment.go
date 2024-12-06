// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0
package types

// VolumeMount type used for mapping pod definition structure.
type VolumeMount struct {
	MountPath string `json:"mountPath,omitempty"`
	Name      string `json:"name,omitempty"`
	Size      string `json:"size,omitempty"`
	Class     string `json:"class,omitempty"`
}

// Container type used for mapping pod definition structure.
type Container struct {
	Name         string        `json:"name,omitempty"`
	Image        string        `json:"image,omitempty"`
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
}

// PersistentVolumeClaim type used for mapping Volume structure.
type PersistentVolumeClaim struct {
	ClaimName string `json:"claimName,omitempty"`
}

// HostPath type used for mapping Volume structure.
type HostPath struct {
	Path string `json:"path,omitempty"`
	Type string `json:"type,omitempty"`
}

// SecretItem for mapping Secret structure.
type SecretItem struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
}

// Secret type for mapping Volume structure.
type Secret struct {
	SecretName string       `json:"secretName,omitempty"`
	Items      []SecretItem `json:"items,omitempty"`
}

// Volume type for mapping Spec structure.
type Volume struct {
	Name                  string                 `json:"name,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaim `json:"persistentVolumeClaim,omitempty"`
	HostPath              *HostPath              `json:"hostPath,omitempty"`
	Secret                *Secret                `json:"secret,omitempty"`
}

// Spec type for mapping Deployment structure.
type Spec struct {
	NodeName      string      `json:"nodeName,omitempty"`
	RestartPolicy string      `json:"restartPolicy,omitempty"`
	Containers    []Container `json:"containers,omitempty"`
	Volumes       []Volume    `json:"volumes,omitempty"`
}

// Deployment type can store k8s deployment data.
type Deployment struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Spec       *Spec  `json:"spec,omitempty"`
}
