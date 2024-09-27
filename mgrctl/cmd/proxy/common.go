// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import "github.com/uyuni-project/uyuni-tools/shared/api"

// Common flag names by key.
const (
	proxyName = "proxyName"
	proxyPort = "proxyPort"
	server    = "server"
	maxCache  = "maxCache"
	email     = "email"
	output    = "output"
)

// Common flags for proxy create config commands.
type ProxyCreateConfigBaseFlags struct {
	ConnectionDetails api.ConnectionDetails `mapstructure:"api"`
	ProxyName         string
	ProxyPort         int
	Server            string
	MaxCache          int
	Email             string
	Output            string
}
