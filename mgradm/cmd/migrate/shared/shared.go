// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// GetSshAuthSocket returns the SSH_AUTH_SOCK environment variable value.
func GetSshAuthSocket() string {
	path := os.Getenv("SSH_AUTH_SOCK")
	if len(path) == 0 {
		log.Fatal().Msg("SSH_AUTH_SOCK is not defined, start an ssh agent and try again")
	}
	return path
}

// GetSshPaths returns the user SSH config and known_hosts paths.
func GetSshPaths() (string, string) {
	// Find ssh config to mount it in the container
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Msg("Failed to find home directory to look for SSH config")
	}
	sshConfigPath := filepath.Join(homedir, ".ssh", "config")
	sshKnownhostsPath := filepath.Join(homedir, ".ssh", "known_hosts")

	if !utils.FileExists(sshConfigPath) {
		sshConfigPath = ""
	}

	if !utils.FileExists(sshKnownhostsPath) {
		sshKnownhostsPath = ""
	}

	return sshConfigPath, sshKnownhostsPath
}
