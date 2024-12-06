// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build !nok8s

package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	adm_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	sshSecretName = "uyuni-migration-key"
	sshConfigName = "uyuni-migration-ssh"
)

func checkSSH(namespace string, flags *adm_utils.SSHFlags) error {
	if exists, err := checkSSHKey(namespace); err != nil {
		return err
	} else if !exists && flags.Key.Public != "" && flags.Key.Private != "" {
		if err := createSSHSecret(namespace, flags.Key.Private, flags.Key.Public); err != nil {
			return err
		}
	} else if !exists {
		return errors.New(L("no SSH key found to use for migration"))
	}

	if exists, err := checkSSHConfig(namespace); err != nil {
		return err
	} else if !exists && flags.Knownhosts != "" {
		// The config may be empty, but not the known_hosts
		if err := createSSHConfig(namespace, flags.Config, flags.Knownhosts); err != nil {
			return err
		}
	} else if !exists {
		return errors.New(L("no SSH known_hosts and configuration found to use for migration"))
	}

	return nil
}

func checkSSHKey(namespace string) (bool, error) {
	exists := false
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "secret", "-n", namespace, sshSecretName, "-o", "jsonpath={.data}",
	)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			log.Debug().Msg("Not found!")
			// The secret was not found, it's not really an error
			return exists, nil
		}
		return exists, utils.Errorf(err, L("failed to get %s SSH key secret"), sshSecretName)
	}
	exists = true

	var data map[string]string
	if err := json.Unmarshal(out, &data); err != nil {
		return exists, err
	}

	for _, key := range []string{"key", "key.pub"} {
		if value, ok := data[key]; !ok || value == "" {
			return exists, fmt.Errorf(L("%[1]s secret misses the %[2]s value"), sshSecretName, key)
		}
	}

	return exists, nil
}

func createSSHSecret(namespace string, keyPath string, pubKeyPath string) error {
	keyContent, err := os.ReadFile(keyPath)
	if err != nil {
		return utils.Errorf(err, L("failed to read key file %s"), keyPath)
	}

	pubContent, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return utils.Errorf(err, L("failed to read public key file %s"), pubKeyPath)
	}

	secret := core.Secret{
		TypeMeta: meta.TypeMeta{APIVersion: "v1", Kind: "Secret"},
		ObjectMeta: meta.ObjectMeta{
			Namespace: namespace,
			Name:      sshSecretName,
			Labels:    kubernetes.GetLabels(kubernetes.ServerApp, ""),
		},
		// It seems serializing this object automatically transforms the secrets to base64.
		Data: map[string][]byte{
			"key":     keyContent,
			"key.pub": pubContent,
		},
	}

	return kubernetes.Apply([]runtime.Object{&secret}, L("failed to create the SSH migration secret"))
}

func checkSSHConfig(namespace string) (bool, error) {
	exists := false
	out, err := utils.RunCmdOutput(
		zerolog.DebugLevel, "kubectl", "get", "cm", "-n", namespace, sshConfigName, "-o", "jsonpath={.data}",
	)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			// The config map was not found, it's not really an error
			return exists, nil
		}
		return exists, utils.Errorf(err, L("failed to get %s SSH ConfigMap"), sshConfigName)
	}
	exists = true

	var data map[string]string
	if err := json.Unmarshal(out, &data); err != nil {
		return exists, utils.Errorf(err, L("failed to parse SSH ConfigMap data"))
	}

	// The known_hosts has to contain at least the entry for the source server.
	if value, ok := data["known_hosts"]; !ok || value == "" {
		return exists, fmt.Errorf(L("%[1]s ConfigMap misses the %[2]s value"), sshSecretName, "known_hosts")
	}

	// An empty config is not an error.
	if _, ok := data["config"]; !ok {
		return exists, fmt.Errorf(L("%[1]s ConfigMap misses the %[2]s value"), sshSecretName, "config")
	}

	return exists, nil
}

func createSSHConfig(namespace string, configPath string, KnownhostsPath string) error {
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return utils.Errorf(err, L("failed to read SSH config file %s"), configPath)
	}

	knownhostsContent, err := os.ReadFile(KnownhostsPath)
	if err != nil {
		return utils.Errorf(err, L("failed to read SSH known_hosts file %s"), KnownhostsPath)
	}

	configMap := core.ConfigMap{
		TypeMeta:   meta.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: meta.ObjectMeta{Namespace: namespace, Name: sshConfigName},
		Data: map[string]string{
			"config":      string(configContent),
			"known_hosts": string(knownhostsContent),
		},
	}
	return kubernetes.Apply([]runtime.Object{&configMap}, L("failed to create the SSH migration ConfigMap"))
}
