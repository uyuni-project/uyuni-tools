// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func kubernetesLogs(
	globalFlags *types.GlobalFlags,
	flags *logsFlags,
	cmd *cobra.Command,
	args []string,
) error {
	cnx := shared.NewConnection("kubectl", "", kubernetes.ProxyFilter)
	podName, err := cnx.GetPodName()
	if err != nil {
		return utils.Errorf(err, L("failed to find proxy pod"))
	}
	namespace, errNamespace := cnx.GetNamespace("")
	if errNamespace != nil {
		return utils.Errorf(err, L("failed to find proxy deployment namespace"))
	}

	commandArgs := []string{"logs", "-n", namespace}
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
		if isRFC3339(flags.Since) {
			commandArgs = append(commandArgs, fmt.Sprintf("--since-time=%s", flags.Since))
		} else {
			commandArgs = append(commandArgs, fmt.Sprintf("--since=%s", flags.Since))
		}
	}

	if len(flags.Containers) == 0 {
		commandArgs = append(commandArgs, podName, "--all-containers")
	} else if len(flags.Containers) == 1 {
		commandArgs = append(commandArgs, flags.Containers[0], "--all-containers")
	} else {
		commandArgs = append(commandArgs, args...)
	}

	return utils.RunCmdStdMapping(zerolog.DebugLevel, "kubectl", commandArgs...)
}

func isRFC3339(timestamp string) bool {
	_, err := time.Parse(time.RFC3339, timestamp)
	return err == nil
}
