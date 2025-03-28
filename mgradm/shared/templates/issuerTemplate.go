// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package templates

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// IssuerTemplate is the base structure for all issuer templates.
type IssuerTemplate struct {
	Namespace string
	FQDN      string
	template  utils.Template
}

// Apply renders the issuer with the certificates and applies them all at once.
func (data IssuerTemplate) Apply() error {
	// Create the server and database certificates
	serverCert, err := CertificateData{
		Namespace:  data.Namespace,
		SecretName: kubernetes.CertSecretName,
		DNSNames:   []string{data.FQDN},
	}.Render()
	if err != nil {
		return err
	}

	dbCert, err := CertificateData{
		Namespace:  data.Namespace,
		SecretName: kubernetes.DBCertSecretName,
		DNSNames:   []string{data.FQDN, "db", "reportdb"},
	}.Render()
	if err != nil {
		return err
	}

	tempDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	path := filepath.Join(tempDir, "issuer.yaml")

	builder := new(strings.Builder)
	if err := data.template.Render(builder); err != nil {
		return utils.Error(err, L("failed to render issuer template"))
	}

	builder.WriteString(serverCert)
	builder.WriteString(dbCert)
	if err := os.WriteFile(path, []byte(builder.String()), 0700); err != nil {
		return utils.Errorf(err, L("failed to write issuer and certificates to %s file"), path)
	}

	_, err = utils.NewRunner("kubectl", "apply", "-f", path).Log(zerolog.DebugLevel).Exec()
	if err != nil {
		return utils.Error(err, L("failed to apply template"))
	}
	return nil
}
