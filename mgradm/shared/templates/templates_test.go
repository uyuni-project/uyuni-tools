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
				Image:        "registry.opensuse.org/uyuni/server-attestation:latest",
				Network:      "uyuni-network",
				DBUserSecret: "db-user-secret",
				DBPassSecret: "db-pass-secret",
			},
		},
		{
			name: "GeneratedCAIssuerTemplate",
			template: NewGeneratedCAIssuerTemplate(
				"default", "uyuni.example.com", "DE", "Bayern", "Nuernberg", "SUSE", "OU", "admin@example.com",
			),
		},
		{
			name: "HubXmlrpcServiceTemplateData",
			template: HubXmlrpcServiceTemplateData{
				CaSecret:   "ca-secret",
				CaPath:     "/etc/pki/ca.crt",
				Ports:      []types.PortMap{utils.NewPortMap("hub", "xmlrpc", 2830, 2830)},
				NamePrefix: "uyuni",
				Image:      "registry.opensuse.org/uyuni/hub-xmlrpc:latest",
				Network:    "uyuni-network",
				ServerHost: "uyuni-server",
			},
		},
		{
			name: "MgrSetupScriptTemplateData",
			template: MgrSetupScriptTemplateData{
				NoSSL:          false,
				DebugJava:      true,
				AdminLogin:     "admin",
				AdminPassword:  "password",
				AdminFirstName: "Administrator",
				AdminLastName:  "Admin",
				AdminEmail:     "admin@example.com",
				OrgName:        "Organization",
			},
		},
		{
			name: "MigrateScriptTemplateData",
			template: MigrateScriptTemplateData{
				Volumes:      []types.VolumeMount{{Name: "var-spacewalk", MountPath: "/var/spacewalk"}},
				SourceFqdn:   "source.example.com",
				User:         "root",
				Kubernetes:   false,
				Prepare:      false,
				DBHost:       "localhost",
				ReportDBHost: "localhost",
			},
		},
		{
			name: "FinalizePostgresTemplateData",
			template: FinalizePostgresTemplateData{
				RunReindex:      true,
				RunSchemaUpdate: true,
				Migration:       true,
				Kubernetes:      false,
			},
		},
		{
			name: "PgsqlMigrateScriptTemplateData",
			template: PgsqlMigrateScriptTemplateData{
				DBHost:       "localhost",
				ReportDBHost: "localhost",
			},
		},
		{
			name: "PgsqlServiceTemplateData",
			template: PgsqlServiceTemplateData{
				Volumes:         []types.VolumeMount{{Name: "var-pgsql", MountPath: "/var/lib/pgsql"}},
				NamePrefix:      "uyuni",
				Ports:           []types.PortMap{utils.NewPortMap("db", "pgsql", 5432, 5432)},
				Image:           "registry.opensuse.org/uyuni/server-postgresql:latest",
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
			name: "PostgreSQLVersionUpgradeTemplateData",
			template: PostgreSQLVersionUpgradeTemplateData{
				OldVersion: "14",
				NewVersion: "16",
				BackupDir:  "/var/lib/pgsql/data-backup",
			},
		},
		{
			name:     "PostUpgradeTemplateData",
			template: PostUpgradeTemplateData{},
		},
		{
			name: "ReusedCAIssuerTemplate",
			template: NewReusedCAIssuerTemplate(
				"default", "uyuni.example.com", "certificate-content", "key-content",
			),
		},
		{
			name: "SalineServiceTemplateData",
			template: SalineServiceTemplateData{
				NamePrefix: "uyuni",
				Image:      "registry.opensuse.org/uyuni/saline:latest",
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
				Ports:           []types.PortMap{utils.NewPortMap("http", "web", 80, 80)},
				Image:           "registry.opensuse.org/uyuni/server:latest",
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
			name: "TLSSecretTemplateData",
			template: TLSSecretTemplateData{
				Name:        "tls-secret",
				Namespace:   "default",
				RootCa:      "root-ca-content",
				Certificate: "cert-content",
				Key:         "key-content",
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

// CertificateData returns string and error, differs from the rest,
// requires separate test case outside the TestTemplatesRender loop.
func TestCertificateDataRender(t *testing.T) {
	data := CertificateData{
		Namespace:  "default",
		SecretName: "uyuni-cert",
		DNSNames:   []string{"uyuni.example.com"},
	}
	if _, err := data.Render(); err != nil {
		t.Errorf("CertificateData render failed: %v", err)
	}
}
