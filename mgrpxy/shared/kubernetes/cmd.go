// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package kubernetes

import (
	"fmt"

	"github.com/spf13/cobra"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// HelmFlags it's used for helm chart flags.
type HelmFlags struct {
	Proxy types.ChartFlags
}

// AddHelmFlags add helm flags to a command.
func AddHelmFlags(cmd *cobra.Command) {
	defaultChart := fmt.Sprintf("oci://%s/%s", utils.DefaultHelmRegistry, utils.DefaultProxyChart)

	cmd.Flags().String("helm-proxy-namespace", "default", L("Kubernetes namespace where to install the proxy"))
	cmd.Flags().String("helm-proxy-chart", defaultChart, L("URL to the proxy helm chart"))
	cmd.Flags().String("helm-proxy-version", "", L("Version of the proxy helm chart"))
	cmd.Flags().String("helm-proxy-values", "", L("Path to a values YAML file to use for proxy helm install"))
}
