// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package squid

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func kubernetesSquidClear(
	globalFlags *types.GlobalFlags,
	flags *squidClearFlags,
	cmd *cobra.Command,
	args []string,
) error {
	cnx := shared.NewKubernetesConnection("", kubernetes.ProxyFilter, "squid")

	if _, err := cnx.Exec("find", "/var/cache/squid", "-mindepth", "1", "-delete"); err != nil {
		return utils.Errorf(err, L("failed to remove cached data"))
	}

	if _, err := cnx.Exec("squid", "-z", "--foreground"); err != nil {
		return utils.Errorf(err, L("failed to re-create the cache directories"))
	}

	return kubernetes.Restart(kubernetes.ProxyFilter)
}
