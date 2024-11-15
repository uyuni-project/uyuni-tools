// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func podmanLogs(
	_ *types.GlobalFlags,
	flags *logsFlags,
	_ *cobra.Command,
	args []string,
) error {
	commandArgs := []string{"logs"}
	if flags.Follow {
		commandArgs = append(commandArgs, "-f")
	}

	if flags.Tail != -1 {
		commandArgs = append(commandArgs, "--tail="+fmt.Sprintf("%d", flags.Tail))
	}

	if flags.Timestamps {
		commandArgs = append(commandArgs, "--timestamps")
	}

	if flags.Since != "" {
		commandArgs = append(commandArgs, fmt.Sprintf("--since=%s", flags.Since))
	}

	if len(flags.Containers) == 0 {
		commandArgs = append(commandArgs, podman.ProxyContainerNames...)
	} else {
		commandArgs = append(commandArgs, args...)
	}

	return utils.RunCmdStdMapping(zerolog.DebugLevel, "podman", commandArgs...)
}
