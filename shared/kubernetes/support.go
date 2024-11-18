// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// RunSupportConfigOnKubernetesHost will run supportconfig command on kubernetes machine.
func RunSupportConfigOnKubernetesHost(dir string, namespace string, filter string) ([]string, error) {
	files, err := utils.RunSupportConfigOnHost()
	if err != nil {
		return files, err
	}

	configmapFilename, err := fetchConfigMap(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any configmap"))
	} else {
		files = append(files, configmapFilename)
	}

	podFilename, err := fetchPodYaml(dir, namespace, filter)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any pod"))
	} else {
		files = append(files, podFilename...)
	}

	return files, nil
}

func fetchConfigMap(dir string, namespace string) (string, error) {
	configmapFile, err := os.Create(path.Join(dir, "configmap"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), configmapFile.Name())
	}
	defer configmapFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "configmap", "-o", "yaml", "-n", namespace)
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch configmap"))
	}

	_, err = configmapFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return configmapFile.Name(), nil
}

func fetchPodYaml(dir string, namespace string, filter string) ([]string, error) {
	pods, err := GetPods(namespace, filter)
	if err != nil {
		return []string{}, utils.Errorf(err, L("cannot check for pods in %s"), filter)
	}

	var podsFile []string
	for _, pod := range pods {
		podFile, err := os.Create(path.Join(dir, fmt.Sprintf("pod-%s", pod)))
		if err != nil {
			log.Warn().Msgf(L("failed to create %s"), podFile.Name())
			continue
		}
		defer podFile.Close()
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", pod, "-o", "yaml", "-n", namespace)
		if err != nil {
			log.Warn().Msgf(L("failed to fetch info for pod %s"), podFile.Name())
			continue
		}

		_, err = podFile.WriteString(string(out))
		if err != nil {
			log.Warn().Msgf(L("failed to write in %s"), podFile.Name())
			continue
		}
		podsFile = append(podsFile, podFile.Name())
	}
	return podsFile, nil
}
