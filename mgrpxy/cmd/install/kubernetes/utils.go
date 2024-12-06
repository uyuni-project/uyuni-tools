// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installForKubernetes(_ *types.GlobalFlags,
	flags *kubernetesProxyInstallFlags, _ *cobra.Command, args []string,
) error {
	for _, binary := range []string{"kubectl", "helm"} {
		if _, err := exec.LookPath(binary); err != nil {
			return fmt.Errorf(L("install %s before running this command"), binary)
		}
	}

	// Unpack the tarball
	configPath := utils.GetConfigPath(args)

	tmpDir, cleaner, err := shared_utils.TempDir()
	if err != nil {
		return err
	}
	defer cleaner()

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
	ports := shared_utils.GetProxyPorts()
	if isK3s {
		err = shared_kubernetes.InstallK3sTraefikConfig(ports)
	} else if IsRke2 {
		err = shared_kubernetes.InstallRke2NginxConfig(ports, flags.Helm.Proxy.Namespace)
	}
	if err != nil {
		return err
	}

	helmArgs := []string{"--set", "ingress=" + clusterInfos.Ingress}
	helmArgs, err = shared_kubernetes.AddSCCSecret(
		helmArgs, flags.Helm.Proxy.Namespace, &flags.SCC, shared_kubernetes.ProxyApp,
	)
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
