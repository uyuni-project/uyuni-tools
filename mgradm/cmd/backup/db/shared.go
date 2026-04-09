// SPDX-FileCopyrightText: 2026 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package db

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func getPostgresConfigPath() (string, error) {
	mountPoint, err := podman.GetVolumeMountPoint(utils.VarPgsqlDataVolumeMount.Name)
	if err != nil {
		return "", err
	}
	return path.Join(mountPoint, "postgresql.conf"), nil
}

// ParsePostgresConfig reads the configuration and returns a map of active settings.
func ParsePostgresConfig() (map[string]string, error) {
	log.Debug().Msg("Reading postgres config")
	configPath, err := getPostgresConfigPath()
	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("Reading %s", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			// Remove inline comments and quotes
			value := strings.Trim(strings.Split(parts[1], "#")[0], " ")
			config[key] = value
		}
	}
	return config, scanner.Err()
}

// UpdatePostgresConfig updates the configuration file with provided key-value pairs.
func UpdatePostgresConfig(updates map[string]string) error {
	log.Debug().Msg("Updating postgres config")
	configPath, err := getPostgresConfigPath()
	if err != nil {
		return err
	}

	log.Trace().Msgf("Writing %s", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	// track which keys we have processed found in the file
	foundKeys := make(map[string]bool)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#") {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) > 0 {
				key := strings.TrimSpace(parts[0])
				if val, ok := updates[key]; ok {
					lines = append(lines, fmt.Sprintf("%s = %s", key, val))
					foundKeys[key] = true
					continue
				}
			}
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// Append missing keys
	for key, val := range updates {
		if !foundKeys[key] {
			lines = append(lines, fmt.Sprintf("%s = %s", key, val))
		}
	}

	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
}
