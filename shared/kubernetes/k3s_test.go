// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// Test that the generated endpoints are valid for traefik.
func Test_GetTraefikEndpointName(t *testing.T) {
	ports := utils.GetServerPorts(true)
	ports = append(ports, utils.HubXmlrpcPorts...)
	ports = append(ports, utils.GetProxyPorts()...)

	for _, port := range ports {
		actual := GetTraefikEndpointName(port)
		// Traefik would fail if the name is longer than 15 characters
		if len(actual) > 15 {
			t.Errorf("Traefik endpoint name has more than 15 characters: %s", actual)
		}
	}
}
