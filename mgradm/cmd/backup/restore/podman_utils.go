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

func parseNetworkData(data []byte) (networkDetails shared.PodanNetworkConfigData, err error) {
	networkData := []map[string]interface{}{}
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
	if _, ok := networkData[0]["network_dns_servers"]; !ok {
		return
	}
	networkDetails.Subnets = networkData[0]["subnets"].([]shared.NetworkSubnet)
	networkDetails.NetworkInsterface = networkData[0]["network_interface"].(string)
	networkDetails.NetworkDNSServers = networkData[0]["network_dns_servers"].([]string)
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
	data := make([]byte, header.Size+1)
	if _, err := tr.Read(data); err != io.EOF {
		log.Warn().Msg(L("Failed to read backed up network configuration, trying default"))
		return errors.Join(err, defaultPodmanNetwork(flags))
	}
	networkDetails, err := parseNetworkData(data)
	if err != nil {
		log.Warn().Msg(L("Failed to decode backed up network configuration, trying default"))
		return errors.Join(err, defaultPodmanNetwork(flags))
	}

	if podman.IsNetworkPresent(podman.UyuniNetwork) {
		if flags.ForceRestore {
			podman.DeleteNetwork(false)
		} else {
			return errors.New(L("podman network already exists"))
		}
	}

	command := []string{"podman", "network", "create", "interface-name", networkDetails.NetworkInsterface}
	for _, v := range networkDetails.Subnets {
		command = append(command, "--subnets", v.Subnet, "--gateway", v.Gateway)
	}
	for _, v := range networkDetails.NetworkDNSServers {
		command = append(command, "--dns", v)
	}
	log.Info().Msgf("Restoring podman network")
	if err := runCmd(command[0], command[1:]...); err != nil {
		log.Error().Msg(L("Unlable to create podman network"))
		return err
	}
	return nil
}

func parseSecretsData(data []byte) (secrets []shared.BackupSecretMap, err error) {
	if err = json.Unmarshal(data, &secrets); err != nil {
		log.Warn().Msg(L("Unable to decode podman secrets"))
		return
	}
	for _, v := range secrets {
		decoded, err := base64.StdEncoding.DecodeString(v.Secret)
		if err != nil {
			log.Warn().Msgf(L("Unable to decode secret %s, using as is"), v.Name)
		} else {
			v.Secret = string(decoded[:])
		}
	}
	return secrets, nil
}

func restorePodmanSecrets(header *tar.Header, tr *tar.Reader, flags *shared.Flagpole) error {
	if flags.DryRun {
		log.Info().Msgf(L("Would restore podman secrets"))
		return nil
	}

	data := make([]byte, header.Size+1)
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
	log.Info().Msgf("Restoring podman secrets")
	baseCommand := []string{"podman", "secret", "create"}
	for _, v := range secrets {
		command := append(baseCommand, v.Name, "-")
		if err := runCmdInput(command[0], v.Secret, command[1:]...); err != nil {
			log.Error().Msg(L("Unlable to create podman secret"))
			hasError = errors.Join(hasError, err)
		}
	}
	return hasError
}
