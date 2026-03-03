// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"os"
	"path"
	"strings"

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

	// Collect cluster-wide information
	clusterInfoFilename, err := fetchClusterInfo(dir)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve cluster info"))
	} else {
		files = append(files, clusterInfoFilename)
	}

	nodeFilename, err := fetchNodeInfo(dir)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve node information"))
	} else {
		files = append(files, nodeFilename)
	}

	eventsFilename, err := fetchEvents(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve events"))
	} else {
		files = append(files, eventsFilename)
	}

	// Collect namespace-specific resources
	configmapFilename, err := fetchConfigMap(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any configmap"))
	} else {
		files = append(files, configmapFilename)
	}

	deploymentFilename, err := fetchDeployments(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve deployments"))
	} else {
		files = append(files, deploymentFilename)
	}

	serviceFilename, err := fetchServices(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve services"))
	} else {
		files = append(files, serviceFilename)
	}

	podFilename, err := fetchPodYaml(dir, namespace, filter)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve any pod"))
	} else {
		files = append(files, podFilename...)
	}

	// Collect Helm release information
	helmFilename, err := fetchHelmReleases(dir, namespace)
	if err != nil {
		log.Warn().Msg(L("cannot retrieve Helm releases"))
	} else {
		files = append(files, helmFilename)
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

func fetchClusterInfo(dir string) (string, error) {
	clusterInfoFile, err := os.Create(path.Join(dir, "cluster-info"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), clusterInfoFile.Name())
	}
	defer clusterInfoFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "cluster-info")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch cluster info"))
	}

	_, err = clusterInfoFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return clusterInfoFile.Name(), nil
}

func fetchNodeInfo(dir string) (string, error) {
	nodeFile, err := os.Create(path.Join(dir, "nodes"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), nodeFile.Name())
	}
	defer nodeFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "nodes", "-o", "yaml")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch node information"))
	}

	_, err = nodeFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return nodeFile.Name(), nil
}

func fetchEvents(dir string, namespace string) (string, error) {
	eventsFile, err := os.Create(path.Join(dir, "events"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), eventsFile.Name())
	}
	defer eventsFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "events", "-n", namespace, "-o", "yaml")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch events"))
	}

	_, err = eventsFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return eventsFile.Name(), nil
}

func fetchDeployments(dir string, namespace string) (string, error) {
	deploymentFile, err := os.Create(path.Join(dir, "deployments"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), deploymentFile.Name())
	}
	defer deploymentFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "deployments", "-n", namespace, "-o", "yaml")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch deployments"))
	}

	_, err = deploymentFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return deploymentFile.Name(), nil
}

func fetchServices(dir string, namespace string) (string, error) {
	serviceFile, err := os.Create(path.Join(dir, "services"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), serviceFile.Name())
	}
	defer serviceFile.Close()
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "kubectl", "get", "services", "-n", namespace, "-o", "yaml")
	if err != nil {
		return "", utils.Errorf(err, L("cannot fetch services"))
	}

	_, err = serviceFile.WriteString(string(out))
	if err != nil {
		return "", err
	}
	return serviceFile.Name(), nil
}

func fetchHelmReleases(dir string, namespace string) (string, error) {
	helmFile, err := os.Create(path.Join(dir, "helm-releases"))
	if err != nil {
		return "", utils.Errorf(err, L("cannot create %s"), helmFile.Name())
	}
	defer helmFile.Close()

	// Get list of Helm releases in the namespace
	out, err := utils.RunCmdOutput(zerolog.DebugLevel, "helm", "list", "-n", namespace, "-a")
	if err != nil {
		return "", utils.Errorf(err, L("cannot list Helm releases"))
	}

	_, err = helmFile.WriteString("==== Helm List ====\n" + string(out) + "\n")
	if err != nil {
		return "", err
	}

	// Get detailed values for each release
	releases := parseHelmReleases(string(out))
	for _, release := range releases {
		out, err := utils.RunCmdOutput(zerolog.DebugLevel, "helm", "get", "values", release, "-n", namespace)
		if err != nil {
			log.Warn().Msgf(L("cannot get values for Helm release %s"), release)
			continue
		}
		_, err = helmFile.WriteString(fmt.Sprintf("==== Helm Values: %s ====\n%s\n", release, string(out)))
		if err != nil {
			log.Warn().Msgf(L("cannot write values for Helm release %s"), release)
			continue
		}
	}

	return helmFile.Name(), nil
}

func parseHelmReleases(helmListOutput string) []string {
	var releases []string
	lines := strings.Split(helmListOutput, "\n")
	// Skip header line and empty lines
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			releases = append(releases, fields[0])
		}
	}
	return releases
}
