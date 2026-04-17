// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"bytes"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/testutils"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestTemplatesRender(t *testing.T) {
	tests := []struct {
		name     string
		template utils.Template
		expected string
	}{
		{
			name: "HttpdTemplateData",
			template: HttpdTemplateData{
				Volumes:        []types.VolumeMount{{Name: "vol1", MountPath: "/mnt/vol1"}},
				HTTPProxyFile:  "/etc/sysconfig/proxy",
				SystemIDSecret: "system-id-secret",
				CaSecret:       "ca-secret",
				CertSecret:     "cert-secret",
				KeySecret:      "key-secret",
			},
		},
		{
			name: "PodTemplateData",
			template: PodTemplateData{
				Ports:         []types.PortMap{{Exposed: 8022, Port: 22}},
				HTTPProxyFile: "/etc/sysconfig/proxy",
				Network:       "uyuni-network",
				IPV6Enabled:   true,
			},
		},
		{
			name: "SaltBrokerTemplateData",
			template: SaltBrokerTemplateData{
				HTTPProxyFile: "/etc/sysconfig/proxy",
			},
		},
		{
			name: "SquidTemplateData",
			template: SquidTemplateData{
				Volumes:       []types.VolumeMount{{Name: "cache-vol", MountPath: "/var/cache/squid"}},
				HTTPProxyFile: "/etc/sysconfig/proxy",
			},
		},
		{
			name: "SSHTemplateData",
			template: SSHTemplateData{
				HTTPProxyFile: "/etc/sysconfig/proxy",
			},
		},
		{
			name: "TFTPDTemplateData",
			template: TFTPDTemplateData{
				Volumes:       []types.VolumeMount{{Name: "tftp-vol", MountPath: "/srv/tftpboot"}},
				HTTPProxyFile: "/etc/sysconfig/proxy",
				CaSecret:      "ca-secret",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tt.template.Render(&buf); err != nil {
				t.Errorf("%s render failed: %v", tt.name, err)
			}
			actual := buf.String()
			if tt.expected != "" && actual != tt.expected {
				diff := testutils.DiffStrings(tt.expected, actual)
				t.Errorf("%s render output mismatch:\n%s", tt.name, diff)
			}
		})
	}
}
