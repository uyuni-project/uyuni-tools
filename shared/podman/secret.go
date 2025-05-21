// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const (
	// DBUserSecret is the name of the podman secret containing the database username.
	DBUserSecret = "uyuni-db-user"
	// DBPassSecret is the name of the podman secret containing the database password.
	DBPassSecret = "uyuni-db-pass"
	// ReportDBUserSecret is the name of the podman secret containing the report database username.
	ReportDBUserSecret = "uyuni-reportdb-user"
	// ReportDBPassSecret is the name of the podman secret containing the report database password.
	ReportDBPassSecret = "uyuni-reportdb-pass"
	// DBUserSecret is the name of the podman secret containing the database admin username.
	DBAdminUserSecret = "uyuni-db-admin-user"
	// DBAdminPassSecret is the name of the podman secret containing the database admin password.
	DBAdminPassSecret = "uyuni-db-admin-pass"
	// CASecret is the name of the podman secret containing the CA certificate.
	CASecret = "uyuni-ca"
	// SSLCertSecret is the name of the podman secret containing the Apache certificate.
	SSLCertSecret = "uyuni-cert"
	// SSLKeySecret is the name of the podman secret containing the Apache SSL certificate key.
	SSLKeySecret = "uyuni-key"
	// DBCASecret is the name of the podman secret containing the Root CA certificate for the database.
	DBCASecret = "uyuni-db-ca"
	// DBSSLCertSecret is the name of the podman secret containing the report database certificate.
	DBSSLCertSecret = "uyuni-db-cert"
	// DBSSLKeySecret is the name of the podman secret containing the report database SSL certificate key.
	DBSSLKeySecret = "uyuni-db-key"
)

// CreateCredentialsSecrets creates the podman secrets, one for the user name and one for the password.
func CreateCredentialsSecrets(userSecret string, user string, passwordSecret string, password string) error {
	if err := createSecret(userSecret, user); err != nil {
		return err
	}
	return createSecret(passwordSecret, password)
}

// CreateCASecrets creates SSL CA.
func CreateCASecrets(
	caSecret string, caPath string,
) error {
	if err := createSecretFromFile(caSecret, caPath); err != nil {
		return utils.Errorf(err, L("failed to create %s secret"), CASecret)
	}
	return nil
}

// CreateTLSSecrets creates SSL CA, Certificate and key secrets.
func CreateTLSSecrets(
	caSecret string, caPath string,
	certSecret string, certPath string,
	keySecret string, keyPath string,
) error {
	if err := createSecretFromFile(certSecret, certPath); err != nil {
		return utils.Errorf(err, L("failed to create %s secret"), SSLCertSecret)
	}

	if err := createSecretFromFile(keySecret, keyPath); err != nil {
		return utils.Errorf(err, L("failed to create %s secret"), SSLKeySecret)
	}
	return CreateCASecrets(caSecret, caPath)
}

// createSecret creates a podman secret.
func createSecret(name string, value string) error {
	tmpDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	secretFile := path.Join(tmpDir, "secret")
	if err := os.WriteFile(secretFile, []byte(value), 0600); err != nil {
		return utils.Errorf(err, L("failed to write %s secret to file"), name)
	}

	return createSecretFromFile(name, secretFile)
}

// createSecretFromFile creates a podman secret from a file.
// Removes any already existing secret with that name.
func createSecretFromFile(name string, secretFile string) error {
	if err := deleteSecret(name, false); err != nil {
		return err
	}

	runner := utils.NewRunner("podman", "secret", "create", name, secretFile).Log(zerolog.DebugLevel)
	if _, err := runner.Exec(); err != nil {
		return utils.Errorf(err, L("failed to create podman secret %s"), name)
	}

	return nil
}

// HasSecret returns whether the secret is defined or not.
func HasSecret(name string) bool {
	return utils.RunCmd("podman", "secret", "exists", name) == nil
}

// DeleteSecret removes a podman secret.
func DeleteSecret(name string, dryRun bool) {
	if err := deleteSecret(name, dryRun); err != nil {
		log.Error().Err(err).Send()
	}
}

func deleteSecret(name string, dryRun bool) error {
	if !HasSecret(name) {
		return nil
	}

	args := []string{"secret", "rm", name}
	command := "podman " + strings.Join(args, " ")
	if dryRun {
		log.Info().Msgf(L("Would run %s"), command)
	} else {
		runner := utils.NewRunner("podman", args...).Log(zerolog.DebugLevel)
		if _, err := runner.Exec(); err != nil {
			return utils.Errorf(err, L("Failed to delete %s secret"), name)
		}
	}
	return nil
}

// GetSecret gets the content of a podman secret given its name.
func GetSecret(name string) (string, error) {
	out, err := newRunner("podman", "secret", "inspect", "--showsecret", name, "--format", "{{.SecretData}}").
		Exec()
	if err != nil {
		return "", utils.Errorf(err, L("failed to get the content of the %s secret"), name)
	}
	return strings.TrimSpace(string(out)), nil
}
