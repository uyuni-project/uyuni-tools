// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// RunSupportConfigOnProxyHost will run supportconfig command on host machine.
func RunSupportConfigOnHost(dir string) ([]string, error) {
	var files []string
	extensions := []string{"", ".md5"}

	// Run supportconfig on the host if installed
	if _, err := exec.LookPath("supportconfig"); err == nil {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "supportconfig")
		if err != nil {
			return []string{}, utils.Errorf(err, L("failed to run supportconfig on the host"))
		}
		tarballPath := utils.GetSupportConfigPath(string(out))

		// Look for the generated supportconfig file
		if tarballPath != "" && utils.FileExists(tarballPath) {
			for _, ext := range extensions {
				files = append(files, tarballPath+ext)
			}
		} else {
			return []string{}, errors.New(L("failed to find host supportconfig tarball from command output"))
		}
	} else {
		log.Warn().Msg(L("supportconfig is not available on the host, skipping it"))
	}

	namespace, err := fetchNamespace(ProxyApp)
	if err != nil {
		return files, err
	}

	configmapFilename, err := fetchConfigMap(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any configmap"))
	} else {
		files = append(files, configmapFilename)
	}

	podFilename, err := fetchPodYaml(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any pod"))
	} else {
		files = append(files, podFilename...)
	}
	return files, nil
}
func fetchNamespace(app string) (string, error) {
	//kubectl get deployment uyuni-proxy -o jsonpath='{.metadata.namespace}'
	namespace, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "deployment", app, "-o=jsonpath={.metadata.namespace}")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch namespace"))
	}
	return string(namespace), nil
}

func fetchConfigMap(dir string, namespace string) (string, error) {
	configmapFile, err := os.Create(path.Join(dir, "configmap"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), configmapFile.Name())
	}
	defer configmapFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "configmap", "-o", "yaml", "--namespace", namespace)
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch configmap"))
	}

	_, err = configmapFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return configmapFile.Name(), nil
}

func fetchPodYaml(dir string, namespace string) ([]string, error) {
	pods, err := GetPods(ProxyFilter)
	if err != nil {
		return []string{}, utils.Errorf(err, L("cannot check for pods in %s"), ProxyFilter)
	}

	var podsFile []string
	for _, pod := range pods {
		podFile, err := os.Create(path.Join(dir, fmt.Sprintf("pod-%s", pod)))
		if err != nil {
			log.Warn().Msgf(L("failed to create %s"), podFile.Name())
			continue
		}
		defer podFile.Close()
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "pod", pod, "-o", "yaml", "--namespace", namespace)
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
