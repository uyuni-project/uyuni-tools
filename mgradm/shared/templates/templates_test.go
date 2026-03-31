// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"io"
	"testing"

	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func TestTemplatesRender(t *testing.T) {
	tests := []struct {
		name     string
		template utils.Template
	}{
		{
			name: "AttestationServiceTemplateData",
			template: AttestationServiceTemplateData{
				NamePrefix:   "uyuni",
				Network:      "uyuni-network",
				DBUserSecret: "db-user-secret",
				DBPassSecret: "db-pass-secret",
			},
		},
		{
			name: "HubXmlrpcServiceTemplateData",
			template: HubXmlrpcServiceTemplateData{
				CaSecret:   "ca-secret",
				CaPath:     "/etc/pki/ca.crt",
				Ports:      []types.PortMap{utils.NewPortMap(2830)},
				NamePrefix: "uyuni",
				Network:    "uyuni-network",
				ServerHost: "uyuni-server",
			},
		},
		{
			name: "PgsqlServiceTemplateData",
			template: PgsqlServiceTemplateData{
				Volumes:         []types.VolumeMount{{Name: "var-pgsql", MountPath: "/var/lib/pgsql"}},
				NamePrefix:      "uyuni",
				Ports:           []types.PortMap{utils.NewPortMap(5432)},
				Network:         "uyuni-network",
				IPV6Enabled:     false,
				CaSecret:        "ca-secret",
				CaPath:          "/etc/pki/ca.crt",
				CertSecret:      "cert-secret",
				CertPath:        "/etc/pki/tls.crt",
				KeySecret:       "key-secret",
				KeyPath:         "/etc/pki/tls.key",
				AdminUser:       "admin-user",
				AdminPassword:   "admin-pass",
				ManagerUser:     "manager-user",
				ManagerPassword: "manager-pass",
				ReportUser:      "report-user",
				ReportPassword:  "report-pass",
			},
		},
		{
			name:     "PostUpgradeTemplateData",
			template: PostUpgradeTemplateData{},
		},
		{
			name: "SalineServiceTemplateData",
			template: SalineServiceTemplateData{
				NamePrefix: "uyuni",
				Network:    "uyuni-network",
				Volumes:    []types.VolumeMount{{Name: "etc-salt", MountPath: "/etc/salt"}},
			},
		},
		{
			name: "PodmanServiceTemplateData",
			template: PodmanServiceTemplateData{
				Volumes:         []types.VolumeMount{{Name: "var-spacewalk", MountPath: "/var/spacewalk"}},
				NamePrefix:      "uyuni",
				Args:            "--arg value",
				Ports:           []types.PortMap{utils.NewPortMap(80)},
				Network:         "uyuni-network",
				IPV6Enabled:     false,
				CaSecret:        "ca-secret",
				CaPath:          "/etc/pki/ca.crt",
				DBCaSecret:      "db-ca-secret",
				DBCaPath:        "/etc/pki/db-ca.crt",
				CertSecret:      "cert-secret",
				CertPath:        "/etc/pki/tls.crt",
				KeySecret:       "key-secret",
				KeyPath:         "/etc/pki/tls.key",
				AdminUser:       "admin-user",
				AdminPassword:   "admin-pass",
				ManagerUser:     "manager-user",
				ManagerPassword: "manager-pass",
				ReportUser:      "report-user",
				ReportPassword:  "report-pass",
			},
		},
		{
			name: "TFTPDTemplateData",
			template: TFTPDTemplateData{
				Network:    "uyuni-network",
				CaSecret:   "ca-secret",
				ServerFQDN: "uyuni.test.org",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.template.Render(io.Discard); err != nil {
				t.Errorf("%s render failed: %v", tt.name, err)
			}
		})
	}
}
