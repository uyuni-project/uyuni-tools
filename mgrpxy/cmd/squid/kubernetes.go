// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package squid

import (
	"github.com/rs/zerolog"

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
	cnx := shared.NewConnection("kubectl", "", kubernetes.ProxyFilter)
	podName, err := cnx.GetPodName()
	if err != nil {
		return utils.Errorf(err, L("failed to get pod name"))
	}

	err = utils.RunCmdStdMapping(zerolog.DebugLevel, "kubectl", "exec", podName, "-c", "squid", "--", "find", "/var/cache/squid", "-mindepth", "1", "-delete")
	if err != nil {
		return utils.Errorf(err, L("failed to remove cached data"))
	}

	return kubernetes.Restart(kubernetes.ProxyFilter)
}
