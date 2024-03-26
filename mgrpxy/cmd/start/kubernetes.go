// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kubernetesStart(
	globalFlags *types.GlobalFlags,
	flags *startFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return kubernetes.Start(kubernetes.ProxyFilter)
}
