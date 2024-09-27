// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// HelmFlags stores Uyuni and Cert Manager Helm information.
type HelmFlags struct {
	Uyuni       types.ChartFlags
	CertManager types.ChartFlags
}

// SslCertFlags can store SSL Certs information.
type SslCertFlags struct {
	Cnames   []string `mapstructure:"cname"`
	Country  string
	State    string
	City     string
	Org      string
	OU       string
	Password string
	Email    string
	Ca       ssl.CaChain
	Server   ssl.SslPair
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
