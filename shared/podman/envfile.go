// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package podman

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// GetServerEnvironmentFile returns the path of the server.env environment file.
func GetServerEnvironmentFile() string {
	return path.Join(GetServiceConfFolder(ServerService), ServerEnvironmentFile)
}

// EnvVarFromEnvironmentFile returns the value of a KEY=VALUE variable of the file, or an empty string.
func EnvVarFromEnvironmentFile(fileName string, name string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Debug().Err(err).Msgf("Failed to read the %s environment file", fileName)
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, name+"=") {
			return strings.TrimSpace(strings.TrimPrefix(line, name+"="))
		}
	}
	return ""
}

// DBFlagsFromEnvironmentFile returns the database connection settings defined in the file.
func DBFlagsFromEnvironmentFile(fileName string) types.DBFlags {
	port, _ := strconv.Atoi(EnvVarFromEnvironmentFile(fileName, "MANAGER_DB_PORT"))
	return types.DBFlags{
		Name:     EnvVarFromEnvironmentFile(fileName, "MANAGER_DB_NAME"),
		Host:     EnvVarFromEnvironmentFile(fileName, "MANAGER_DB_HOST"),
		Port:     port,
		Provider: EnvVarFromEnvironmentFile(fileName, "EXTERNALDB_PROVIDER"),
	}
}
