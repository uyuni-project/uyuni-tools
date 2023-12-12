// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type HelmFlags struct {
	Proxy types.ChartFlags
}

func AddHelmFlags(cmd *cobra.Command) {
	defaultChart := fmt.Sprintf("oci://%s/proxy-helm", utils.DefaultNamespace)

	cmd.Flags().String("helm-proxy-namespace", "default", "Kubernetes namespace where to install the proxy")
	cmd.Flags().String("helm-proxy-chart", defaultChart, "URL to the proxy helm chart")
	cmd.Flags().String("helm-proxy-version", "", "Version of the proxy helm chart")
	cmd.Flags().String("helm-proxy-values", "", "Path to a values YAML file to use for proxy helm install")
}
