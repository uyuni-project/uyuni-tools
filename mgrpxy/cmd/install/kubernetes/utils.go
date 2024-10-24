// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installForKubernetes(globalFlags *types.GlobalFlags,
	flags *kubernetesProxyInstallFlags, cmd *cobra.Command, args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf(L("install %s before running this command"), binary)
		}
	}

	// Unpack the tarball
	configPath := utils.GetConfigPath(args)

	tmpDir, err := shared_utils.TempDir()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := shared_utils.ExtractTarGz(configPath, tmpDir); err != nil {
		return errors.New(L("failed to extract configuration"))
	}

	// Check the kubernetes cluster setup
	clusterInfos, err := shared_kubernetes.CheckCluster()
	if err != nil {
		return err
	}

	// If installing on k3s, install the traefik helm config in manifests
	isK3s := clusterInfos.IsK3s()
	IsRke2 := clusterInfos.IsRke2()
	if isK3s {
		shared_kubernetes.InstallK3sTraefikConfig(shared_utils.ProxyTCPPorts, shared_utils.UDPPorts)
	} else if IsRke2 {
		shared_kubernetes.InstallRke2NginxConfig(shared_utils.ProxyTCPPorts, shared_utils.UDPPorts,
			flags.Helm.Proxy.Namespace)
	}

	helmArgs := []string{"--set", "ingress=" + clusterInfos.Ingress}
	helmArgs, err = shared_kubernetes.AddSccSecret(helmArgs, flags.Helm.Proxy.Namespace, &flags.Scc)
	if err != nil {
		return err
	}

	// Install the uyuni proxy helm chart
	if err := kubernetes.Deploy(
		&flags.ProxyImageFlags, &flags.Helm, tmpDir, clusterInfos.GetKubeconfig(), helmArgs...,
	); err != nil {
		return shared_utils.Errorf(err, L("cannot deploy proxy helm chart"))
	}

	return nil
}
