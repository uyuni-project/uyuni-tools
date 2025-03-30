// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

const (
	// DBUserSecret is the name of the podman secret containing the database username.
	DBUserSecret = "uyuni-db-user"
	// DBUserSecret is the name of the podman secret containing the database password.
	DBPassSecret = "uyuni-db-pass"
)

// CreateDBSecrets creates the podman secrets for the database credentials.
func CreateDBSecrets(user string, password string) error {
	if err := createSecret(DBUserSecret, user); err != nil {
		return err
	}
	return createSecret(DBPassSecret, password)
}

// createSecret creates a podman secret.
func createSecret(name string, value string) error {
	if hasSecret(name) {
		return nil
	}

	tmpDir, cleaner, err := utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	secretFile := path.Join(tmpDir, "secret")
	if err := os.WriteFile(secretFile, []byte(value), 0o600); err != nil {
		return utils.Errorf(err, L("failed to write %s secret to file"), name)
	}

	if err := utils.RunCmd("podman", "secret", "create", name, secretFile); err != nil {
		return utils.Errorf(err, L("failed to create podman secret %s"), name)
	}

	return nil
}

func hasSecret(name string) bool {
	return utils.RunCmd("podman", "secret", "exists", name) == nil
}

// DeleteSecret removes a podman secret.
func DeleteSecret(name string, dryRun bool) {
	if !hasSecret(name) {
		return
	}

	args := []string{"secret", "rm", name}
	command := "podman " + strings.Join(args, " ")
	if dryRun {
		log.Info().Msgf(L("Would run %s"), command)
	} else {
		log.Info().Msgf(L("Run %s"), command)
		if err := utils.RunCmd("podman", args...); err != nil {
			log.Error().Err(err).Msgf(L("Failed to delete %s secret"), name)
		}
	}
}
