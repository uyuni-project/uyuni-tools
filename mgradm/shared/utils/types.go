// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/ssl"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// HelmFrags stores Uyuni and Cert Manager Helm information.
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
