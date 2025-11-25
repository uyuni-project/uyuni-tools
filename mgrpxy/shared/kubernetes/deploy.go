// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

const helmAppName = "uyuni-proxy"

// KubernetesProxyUpgradeFlags represents the flags for the mgrpxy upgrade kubernetes command.
type KubernetesProxyUpgradeFlags struct {
	utils.ProxyImageFlags `mapstructure:",squash"`
	SCC                   types.SCCCredentials
	Helm                  HelmFlags
}

// Deploy will deploy proxy in kubernetes.
func Deploy(imageFlags *utils.ProxyImageFlags, helmFlags *HelmFlags, configDir string,
	kubeconfig string, helmArgs ...string,
) error {
	log.Info().Msg(L("Installing Uyuni proxy"))

	helmParams := []string{}

	// Pass the user-provided values file
	extraValues := helmFlags.Proxy.Values
	if extraValues != "" {
		helmParams = append(helmParams, "-f", extraValues)
	}

	if !shared_utils.FileExists(path.Join(configDir, "httpd.yaml")) {
		if _, err := getHTTPDYaml(configDir); err != nil {
			return err
		}
	}
	helmParams = append(helmParams, "-f", path.Join(configDir, "httpd.yaml"))

	if !shared_utils.FileExists(path.Join(configDir, "ssh.yaml")) {
		if _, err := getSSHYaml(configDir); err != nil {
			return err
		}
	}
	helmParams = append(helmParams, "-f", path.Join(configDir, "ssh.yaml"))

	if !shared_utils.FileExists(path.Join(configDir, "config.yaml")) {
		if _, err := getConfigYaml(configDir); err != nil {
			return err
		}
	}
	helmParams = append(helmParams, "-f", path.Join(configDir, "config.yaml"))

	if len(imageFlags.Tuning.Httpd) > 0 {
		absPath, err := filepath.Abs(imageFlags.Tuning.Httpd)
		if err != nil {
			return err
		}
		helmParams = append(helmParams, "--set-file", "apache_tuning="+absPath)
	}

	if len(imageFlags.Tuning.Squid) > 0 {
		absPath, err := filepath.Abs(imageFlags.Tuning.Squid)
		if err != nil {
			return err
		}
		helmParams = append(helmParams, "--set-file", "squid_tuning="+absPath)
	}

	if len(imageFlags.Tuning.SSH) > 0 {
		absPath, err := filepath.Abs(imageFlags.Tuning.SSH)
		if err != nil {
			return err
		}
		helmParams = append(helmParams, "--set-file", "ssh_tuning="+absPath)
	}

	helmParams = append(helmParams,
		"--set", "images.proxy-httpd="+imageFlags.GetContainerImage("httpd"),
		"--set", "images.proxy-salt-broker="+imageFlags.GetContainerImage("salt-broker"),
		"--set", "images.proxy-squid="+imageFlags.GetContainerImage("squid"),
		"--set", "images.proxy-ssh="+imageFlags.GetContainerImage("ssh"),
		"--set", "images.proxy-tftpd="+imageFlags.GetContainerImage("tftpd"),
		"--set", "repository="+imageFlags.Registry.Host,
		"--set", "version="+imageFlags.Tag,
		"--set", "pullPolicy="+string(kubernetes.GetPullPolicy(imageFlags.PullPolicy)))

	helmParams = append(helmParams, helmArgs...)

	// Install the helm chart
	if err := kubernetes.HelmUpgrade(kubeconfig, helmFlags.Proxy.Namespace, true, "", helmAppName, helmFlags.Proxy.Chart,
		helmFlags.Proxy.Version, helmParams...); err != nil {
		return shared_utils.Errorf(err, L("cannot run helm upgrade"))
	}

	// Wait for the pod to be started
	return kubernetes.WaitForDeployments(helmFlags.Proxy.Namespace, helmAppName)
}

func getSSHYaml(directory string) (string, error) {
	sshPayload, err := kubernetes.GetSecret("proxy-secret", "-o=jsonpath={.data.ssh\\.yaml}")
	if err != nil {
		return "", err
	}

	sshYamlFilename := filepath.Join(directory, "ssh.yaml")
	err = os.WriteFile(sshYamlFilename, []byte(sshPayload), 0644)
	if err != nil {
		return "", shared_utils.Errorf(err, L("failed to write in file %s"), sshYamlFilename)
	}

	return sshYamlFilename, nil
}

func getHTTPDYaml(directory string) (string, error) {
	httpdPayload, err := kubernetes.GetSecret("proxy-secret", "-o=jsonpath={.data.httpd\\.yaml}")
	if err != nil {
		return "", err
	}

	httpdYamlFilename := filepath.Join(directory, "httpd.yaml")
	err = os.WriteFile(httpdYamlFilename, []byte(httpdPayload), 0644)
	if err != nil {
		return "", shared_utils.Errorf(err, L("failed to write in file %s"), httpdYamlFilename)
	}

	return httpdYamlFilename, nil
}

func getConfigYaml(directory string) (string, error) {
	configPayload, err := kubernetes.GetConfigMap("proxy-configMap", "-o=jsonpath={.data.config\\.yaml}")
	if err != nil {
		return "", err
	}

	configYamlFilename := filepath.Join(directory, "config.yaml")
	err = os.WriteFile(configYamlFilename, []byte(configPayload), 0644)
	if err != nil {
		return "", shared_utils.Errorf(err, L("failed to write in file %s"), configYamlFilename)
	}

	return configYamlFilename, nil
}

// Upgrade will upgrade the current kubernetes proxy.
func Upgrade(flags *KubernetesProxyUpgradeFlags, _ *cobra.Command, _ []string) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf(L("install %s before running this command"), binary)
		}
	}

	tmpDir, cleaner, err := shared_utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

	// Check the kubernetes cluster setup
	clusterInfos, err := kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	namespace := flags.Helm.Proxy.Namespace
	if _, err = kubernetes.GetNode(namespace, kubernetes.ProxyApp); err != nil {
		if err := kubernetes.ReplicasTo(namespace, kubernetes.ProxyApp, 1); err != nil {
			return err
		}
	}

	err = kubernetes.ReplicasTo(namespace, kubernetes.ProxyApp, 0)
	if err != nil {
		return err
	}

	defer func() {
		// if something is running, we don't need to set replicas to 1
		if _, err = kubernetes.GetNode(namespace, kubernetes.ProxyApp); err != nil {
			if err = kubernetes.ReplicasTo(namespace, kubernetes.ProxyApp, 1); err != nil {
				log.Error().Err(err).Msg(L("failed to scale replicas to 1"))
			}
		}
	}()

	helmArgs := []string{"--set", "ingress=" + clusterInfos.Ingress}

	// Get the registry secret name if any
	pullSecret, err := kubernetes.GetDeploymentImagePullSecret(namespace, kubernetes.ProxyFilter)
	if err != nil {
		return err
	}
	if pullSecret != "" {
		helmArgs = append(helmArgs, "--set", "registrySecret="+pullSecret)
	}

	// Install the uyuni proxy helm chart
	if err := Deploy(&flags.ProxyImageFlags, &flags.Helm, tmpDir, clusterInfos.GetKubeconfig(),
		helmArgs...,
	); err != nil {
		return shared_utils.Errorf(err, L("cannot deploy proxy helm chart"))
	}

	return nil
}
