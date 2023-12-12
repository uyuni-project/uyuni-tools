// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/mgrpxy/shared/utils"
	shared_kubernetes "github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

func installForKubernetes(globalFlags *types.GlobalFlags,
	flags *kubernetesProxyInstallFlags, cmd *cobra.Command, args []string,
) {
	// Unpack the tarball
	configPath := utils.GetConfigPath(args)

	tmpDir, err := os.MkdirTemp("", "mgrpxy-*")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create temporary directory")
	}
	defer os.RemoveAll(tmpDir)

	if err := shared_utils.ExtractTarGz(configPath, tmpDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to extract configuration")
	}

	// Check the kubernetes cluster setup
	clusterInfos := shared_kubernetes.CheckCluster()

	// If installing on k3s, install the traefik helm config in manifests
	isK3s := clusterInfos.IsK3s()
	IsRke2 := clusterInfos.IsRke2()
	if isK3s {
		shared_kubernetes.InstallK3sTraefikConfig(shared_utils.PROXY_TCP_PORTS, shared_utils.UDP_PORTS)
	} else if IsRke2 {
		shared_kubernetes.InstallRke2NginxConfig(shared_utils.PROXY_TCP_PORTS, shared_utils.UDP_PORTS,
			flags.Helm.Proxy.Namespace)
	}

	// Install the uyuni proxy helm chart
	kubernetes.Deploy(&flags.ProxyInstallFlags, &flags.Helm, tmpDir, clusterInfos.GetKubeconfig(),
		"--set", "ingress="+clusterInfos.Ingress)
}
