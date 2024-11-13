// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// InstallSSLFlags holds all the flags values related to SSL for installation.
type InstallSSLFlags struct {
	types.SSLCertGenerationFlags `mapstructure:",squash"`
	Ca                           types.CaChain
	Server                       types.SSLPair
}

// KubernetesFlags stores Uyuni and Cert Manager kubernetes specific parameters.
type KubernetesFlags struct {
	Uyuni       types.ChartFlags
	CertManager types.ChartFlags
}

// HubXmlrpcFlags contains settings for Hub XMLRPC container.
type HubXmlrpcFlags struct {
	Replicas  int
	Image     types.ImageFlags `mapstructure:",squash"`
	IsChanged bool
}

// CocoFlags contains settings for coco attestation container.
type CocoFlags struct {
	Replicas  int
	Image     types.ImageFlags `mapstructure:",squash"`
	IsChanged bool
}

// SalineFlags contains settings for Saline container.
type SalineFlags struct {
	Port      int
	Replicas  int
	Image     types.ImageFlags `mapstructure:",squash"`
	IsChanged bool
}

// VolumeFlags stores the persistent volume claims configuration.
type VolumesFlags struct {
	// Class is the default storage class for all the persistent volume claims.
	Class string
	// Database is the configuration of the var-pgsql volume.
	Database VolumeFlags
	// Packages is the configuration of the var-spacewalk volume containing the synchronizede repositories.
	Packages VolumeFlags
	// Www is the configuration of the srv-www volume containing the imags and distributions.
	Www VolumeFlags
	// Cache is the configuration of the var-cache volume.
	Cache VolumeFlags
	// Mirror is the PersistentVolume name to use in case of a mirror setup.
	// An empty value means no mirror will be used.
	Mirror string
}

// VolumeFlags is the configuration of one volume.
type VolumeFlags struct {
	// Size is the requested size of the volume using kubernetes values like '100Gi'.
	Size string
	// Class is the storage class of the volume.
	Class string
}
