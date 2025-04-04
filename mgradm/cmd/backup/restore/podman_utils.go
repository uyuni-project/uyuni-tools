// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package restore

import (
	"archive/tar"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
)

func restorePodmanConfiguration(podmanBackupFile string, flags *shared.Flagpole) error {
	// Read tarball
	backupFile, err := os.Open(podmanBackupFile)
	if err != nil {
		return err
	}
	defer backupFile.Close()

	var hasError error

	tr := tar.NewReader(backupFile)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch header.Name {
		case shared.NetworkOutputFile:
			hasError = errors.Join(hasError, restorePodmanNetwork(header, tr, flags))
		case shared.SecretBackupFile:
			hasError = errors.Join(hasError, restorePodmanSecrets(header, tr, flags))
		default:
			log.Warn().Msgf(L("Ignoring unexpected file in the podman backup %s"), header.Name)
		}
	}
	return hasError
}

// parseNetworkData decodes stored podman network inspect result.
// We are not interested in all data, so selectively decode intereting bits.
func parseNetworkData(data []byte) (networkDetails shared.PodanNetworkConfigData, err error) {
	var networkData []map[string]json.RawMessage
	if err = json.Unmarshal(data, &networkData); err != nil {
		log.Warn().Msg(L("Unable to decode network data backup"))
		return
	}

	err = errors.New(L("Incorrect network data backup"))
	if len(networkData) != 1 {
		return
	}
	if _, ok := networkData[0]["subnets"]; !ok {
		return
	}
	if _, ok := networkData[0]["network_interface"]; !ok {
		return
	}

	// Optional
	if _, ok := networkData[0]["network_dns_servers"]; ok {
		if err = json.Unmarshal(networkData[0]["network_dns_servers"], &networkDetails.NetworkDNSServers); err != nil {
			return
		}
	}

	if err = json.Unmarshal(networkData[0]["subnets"], &networkDetails.Subnets); err != nil {
		return
	}

	if err = json.Unmarshal(networkData[0]["network_interface"], &networkDetails.NetworkInsterface); err != nil {
		return
	}

	return networkDetails, nil
}

func defaultPodmanNetwork(flags *shared.Flagpole) error {
	if podman.IsNetworkPresent(podman.UyuniNetwork) {
		if flags.ForceRestore {
			podman.DeleteNetwork(false)
		} else {
			return errors.New(L("podman network already exists"))
		}
	}
	if err := podman.SetupNetwork(false); err != nil {
		log.Error().
			Msg(L("Unable to create podman network! Check the error and create network manually before starting the service"))
		return err
	}
	return nil
}

func restorePodmanNetwork(header *tar.Header, tr *tar.Reader, flags *shared.Flagpole) error {
	if flags.DryRun {
		log.Info().Msgf(L("Would restore network configuration"))
		return nil
	}
	data := make([]byte, header.Size)
	if _, err := tr.Read(data); err != io.EOF {
		log.Warn().Msg(L("Failed to read backed up network configuration, trying default"))
		return errors.Join(err, defaultPodmanNetwork(flags))
	}
	log.Trace().Msgf("Loaded network data: %s", data)
	networkDetails, err := parseNetworkData(data)
	if err != nil {
		log.Warn().Err(err).Msg(L("Failed to decode backed up network configuration, trying default"))
		return errors.Join(err, defaultPodmanNetwork(flags))
	}

	if podman.IsNetworkPresent(podman.UyuniNetwork) {
		if flags.ForceRestore {
			podman.DeleteNetwork(false)
		} else {
			log.Warn().Msg(L("Podman network already exists, not restoring unless forced"))
			return errors.New(L("podman network already exists"))
		}
	}

	command := []string{"podman", "network", "create", "--interface-name", networkDetails.NetworkInsterface}
	for _, v := range networkDetails.Subnets {
		command = append(command, "--subnet", v.Subnet, "--gateway", v.Gateway)
	}
	for _, v := range networkDetails.NetworkDNSServers {
		command = append(command, "--dns", v)
	}
	command = append(command, podman.UyuniNetwork)
	log.Info().Msg(L("Restoring podman network"))
	if err := runCmd(command[0], command[1:]...); err != nil {
		log.Error().Err(err).Msg(L("Unlable to create podman network"))
		return err
	}
	return nil
}

func parseSecretsData(data []byte) ([]shared.BackupSecretMap, error) {
	secrets := []shared.BackupSecretMap{}
	if err := json.Unmarshal(data, &secrets); err != nil {
		log.Warn().Err(err).Msg(L("Unable to decode podman secrets"))
		return nil, err
	}

	decodedSecrets := make([]shared.BackupSecretMap, len(secrets))
	for i, v := range secrets {
		decoded, err := base64.StdEncoding.DecodeString(v.Secret)
		if err != nil {
			log.Warn().Msgf(L("Unable to decode secret %s, using as is"), v.Name)
		} else {
			decodedSecrets[i] = shared.BackupSecretMap{
				Name:   v.Name,
				Secret: string(decoded[:]),
			}
		}
	}
	return decodedSecrets, nil
}

func restorePodmanSecrets(header *tar.Header, tr *tar.Reader, flags *shared.Flagpole) error {
	if flags.DryRun {
		log.Info().Msgf(L("Would restore podman secrets"))
		return nil
	}

	data := make([]byte, header.Size)
	if _, err := tr.Read(data); err != io.EOF {
		log.Warn().Msg(L("Failed to read backed up podman secrets, no secrets were restored"))
		return err
	}
	secrets, err := parseSecretsData(data)
	if err != nil {
		log.Warn().Msg(L("Failed to decode backed up podman secrets, no secrets were restored"))
		return err
	}

	var hasError error
	log.Info().Msg(L("Restoring podman secrets"))
	baseCommand := []string{"podman", "secret", "create"}
	for _, v := range secrets {
		if podman.IsSecretPresent(v.Name) {
			if flags.ForceRestore {
				baseCommand = append(baseCommand, "--replace")
			} else {
				log.Error().Msgf(L("Podman secret %s is already present, not restoring unless forced"), v.Name)
				continue
			}
		}
		command := append(baseCommand, v.Name, "-")
		if err := runCmdInput(command[0], v.Secret, command[1:]...); err != nil {
			log.Error().Msg(L("Unable to create podman secret"))
			hasError = errors.Join(hasError, err)
		}
	}
	return hasError
}

// handleVolumeHacks handles special import cicrumstances.
// etc-apache2 volume has a trailing mime.type link which podman import reports as an error.
// etc-postfix volume has a trailing cacerts link which podman import reports as an error.
func handleVolumeHacks(volume string, inErr error) error {
	// Quick return if everything is ok
	// or if something happened and volume was not imported at all
	if inErr == nil || !podman.IsVolumePresent(volume) {
		return inErr
	}
	// special handling of apache2 and postfix volume with trailing links
	basePath, err := podman.GetPodmanVolumeBasePath()
	if err != nil {
		log.Debug().Msg("cannot get base volume path")
		return inErr
	}
	switch volume {
	case "etc-apache2":
		log.Debug().Msg("Special apache2 volume handling")
		volumePath := path.Join(basePath, "etc-apache2", "_data")
		if err := os.Symlink("../mime.types", path.Join(volumePath, "mime.types")); err != nil {
			log.Debug().Err(err).Msgf("cannot create link %s for %s", path.Join(volumePath, "mime.types"), "../mime.types")
			return inErr
		}
		return nil
	case "etc-postfix":
		log.Debug().Msg("Special postfix volume handling")
		volumePath := path.Join(basePath, "etc-postfix", "_data")
		if err := os.Symlink("../../ssl/certs", path.Join(volumePath, "ssl", "cacerts")); err != nil {
			log.Debug().Err(err).Msgf("cannot create link %s for %s", path.Join(volumePath, "ssl", "cacerts"), "../../ssl/certs")
			return inErr
		}
		return nil
	}
	return inErr
}
