// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package proxy

import (
	"strings"
)

// GetFilename returns the filename to save a configuration file to.
// If an output filename is not specified, then the filename is based on the proxy name.
func GetFilename(output string, proxyName string) string {
	filename := output
	if filename == "" {
		filename = strings.Split(proxyName, ".")[0] + "-config"
	}
	return filename + ".tar.gz"
}
