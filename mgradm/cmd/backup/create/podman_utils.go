// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	backup "github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// ExportConfiguration creates a tarball with podman network and secrets configuration.
// If dryRun is true, only messages explaining what to do will be logged.
func exportPodmanConfiguration(outputDir string, dryRun bool) error {
	network, errNetwork := backupPodmanNetwork(dryRun)
	secrets, errPodman := backupPodmanSecrets(dryRun)

	if dryRun {
		return nil
	}
	// Create output file
	out, err := os.Create(path.Join(outputDir, backup.PodmanConfBackupFile))
	if err != nil {
		return fmt.Errorf(L("failed to create podman backup tarball: %w"), err)
	}
	defer out.Close()

	// Prepare tar buffer
	tw := tar.NewWriter(out)
	defer tw.Close()
	var hasError error

	if errNetwork != nil {
		log.Warn().Msg(L("Network was not backed up"))
	} else {
		header := &tar.Header{
			Name: backup.NetworkOutputFile,
			Mode: 0622,
			Size: int64(len(network)),
		}
		if err := tw.WriteHeader(header); err != nil {
			hasError = err
		}
		if _, err := tw.Write(network); err != nil {
			hasError = utils.JoinErrors(hasError, err)
		}
	}

	if errPodman != nil {
		log.Warn().Msg(L("Podman secrets were not backed up"))
	} else {
		header := &tar.Header{
			Name: backup.SecretBackupFile,
			Mode: 0600,
			Size: int64(len(secrets)),
		}
		if err := tw.WriteHeader(header); err != nil {
			hasError = utils.JoinErrors(hasError, err)
		}
		if _, err := tw.Write(secrets); err != nil {
			hasError = utils.JoinErrors(hasError, err)
		}
	}
	return hasError
}

func backupPodmanNetwork(dryRun bool) ([]byte, error) {
	networkExportCommand := []string{"podman", "network", "inspect", podman.UyuniNetwork}
	if dryRun {
		log.Info().Msgf(L("Would run %s"), strings.Join(networkExportCommand, " "))
		return nil, nil
	}
	output, err := runCmdOutput(zerolog.DebugLevel, networkExportCommand[0], networkExportCommand[1:]...)
	if err != nil {
		log.Warn().Err(err).Msg(L("Failed to export network data"))
		return nil, err
	}
	return output, nil
}

func backupPodmanSecrets(dryRun bool) ([]byte, error) {
	const secretFile = "/var/lib/containers/storage/secrets/filedriver/secretsdata.json"
	secretListCommand := []string{"podman", "secret", "ls", "--format", "{{range .}}{{.Name}}:{{.ID}},{{end}}"}
	if dryRun {
		log.Info().Msgf(L("Would run %s"), strings.Join(secretListCommand, " "))
		return nil, nil
	}
	output, err := runCmdOutput(zerolog.DebugLevel, secretListCommand[0], secretListCommand[1:]...)
	if err != nil {
		log.Warn().Err(err).Msg(L("Failed to export secrets data"))
		return nil, err
	}

	type SecretMap struct {
		Name string
		ID   string
	}
	secretMappings := []SecretMap{}
	for _, v := range strings.Split(string(output), ",") {
		tmp := strings.SplitN(v, ":", 2)
		// Ignore different length, usually last emptry string
		if len(tmp) == 2 {
			secretMappings = append(secretMappings, SecretMap{Name: tmp[0], ID: tmp[1]})
		}
	}

	// load secretFile as json
	output, err = os.ReadFile(secretFile)
	if err != nil {
		log.Warn().Err(err).Msg(L("Failed to read secrets data"))
		return nil, err
	}

	var podmanSecrets map[string]string
	if err := json.Unmarshal(output, &podmanSecrets); err != nil {
		log.Warn().Err(err).Msg(L("Unable to decode podman secrets"))
		return nil, err
	}

	// store id -> secret file in tar ball location
	backupSecretMap := []backup.BackupSecretMap{}
	for _, secretMap := range secretMappings {
		for secretID, secretValue := range podmanSecrets {
			if secretMap.ID == secretID {
				backupSecretMap = append(backupSecretMap, backup.BackupSecretMap{Name: secretMap.Name, Secret: secretValue})
			}
		}
	}
	output, err = json.Marshal(backupSecretMap)
	if err != nil {
		log.Warn().Err(err).Msg(L("Unable to encode secrets backup"))
		return nil, err
	}
	return output, nil
}
