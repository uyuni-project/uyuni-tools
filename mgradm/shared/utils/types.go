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

// HelmFlags stores Uyuni and Cert Manager Helm information.
type HelmFlags struct {
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
