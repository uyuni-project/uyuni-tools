// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared"
	"github.com/uyuni-project/uyuni-tools/shared/kubernetes"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type logsFlags struct {
	Containers []string
	Follow     bool
	Timestamps bool
	Tail       int
	Since      string
}

var systemd podman.Systemd = podman.SystemdImpl{}

// NewCommand to get the logs of the server.
func newCmd(globalFlags *types.GlobalFlags, run utils.CommandFunc[logsFlags]) *cobra.Command {
	var flags logsFlags

	cmd := &cobra.Command{
		Use:   "logs [pod [container] | container...]",
		Short: L("Get the proxy logs"),
		Long: L(`Get the proxy logs
The command automatically detects installed backend and displays the logs for containers managed by Kubernetes or Podman
However, you can specify the pod and/or container names to get the logs for specific container(s).
See examples for more details.`),
		Example: `  Log all relevant containers (Podman and Kubernetes)

    $ mgrpxy logs

  Log all relevant containers in the specified pod (Kubernetes)

    $ mgrpxy logs uyuni-proxy-pod

  Log the specified container in the specified pod (Kubernetes)

    $ mgrpxy logs uyuni-proxy-pod httpd

  Log the specified containers (Podman)

    $ mgrpxy logs logs uyuni-proxy-httpd uyuni-proxy-ssh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.Containers = cmd.Flags().Args()
			return utils.CommandHelper(globalFlags, cmd, args, &flags, nil, run)
		},
		ValidArgsFunction: getContainerNames,
	}

	cmd.Flags().BoolP("follow", "f", false, L("specify if logs should be followed"))
	cmd.Flags().BoolP("timestamps", "t", false, L("show timestamps in the log outputs"))
	cmd.Flags().Int("tail", -1, L("number of lines to show from the end of the logs"))
	cmd.Flags().Lookup("tail").NoOptDefVal = "-1"
	cmd.Flags().String("since", "",
		L(`show logs since a specific time or duration.
Supports Go duration strings and RFC3339 format (e.g. 3h, 2023-01-02T15:04:05)`),
	)

	cmd.SetUsageTemplate(cmd.UsageTemplate())
	return cmd
}

// NewCommand to get the logs of the server.
func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	return newCmd(globalFlags, logs)
}

func logs(globalFlags *types.GlobalFlags, flags *logsFlags, cmd *cobra.Command, args []string) error {
	fn, err := shared.ChooseProxyPodmanOrKubernetes(cmd.Flags(), podmanLogs, kubernetesLogs)
	if err != nil {
		return err
	}

	return fn(globalFlags, flags, cmd, args)
}

func getContainerNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	var names []string

	if systemd.HasService(podman.ProxyService) {
		names = getNames(exec.Command("podman", "ps", "--format", "{{.Names}}"), "\n", "uyuni")
	} else if utils.IsInstalled("kubectl") && utils.IsInstalled("helm") {
		if len(args) == 0 {
			cnx := shared.NewConnection("kubectl", "", kubernetes.ProxyFilter)
			podName, err := cnx.GetPodName()
			if err != nil {
				log.Fatal().Err(err)
			}
			return []string{podName}, cobra.ShellCompDirectiveNoFileComp
		} else if len(args) == 1 {
			names = getNames(
				exec.Command("kubectl", "get", "pod", args[0], "-o", "jsonpath={.spec.containers[*].name}"),
				" ", "",
			)
		} else {
			// kubernetes log only accepts either 1 container name or the --all-containers flag.
			return names, cobra.ShellCompDirectiveNoFileComp
		}
	}

	return minus(names, args), cobra.ShellCompDirectiveNoFileComp
}

// retrieves pod/container retrieve command and parses its names for auto completion.
func getNames(cmd *exec.Cmd, cmdResultSeparator string, namesPrefix string) []string {
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	names := strings.Split(strings.TrimSpace(string(out)), cmdResultSeparator)
	if namesPrefix == "" {
		return names
	}

	var filteredNames []string
	for _, name := range names {
		if strings.HasPrefix(name, namesPrefix) {
			filteredNames = append(filteredNames, name)
		}
	}
	return filteredNames
}

// Returns the elements of a left slice minus the elements of the right slice.
func minus(left []string, right []string) []string {
	rightMap := make(map[string]bool)
	for _, elementRight := range right {
		rightMap[elementRight] = true
	}

	var result []string
	for _, elementLeft := range left {
		if !rightMap[elementLeft] {
			result = append(result, elementLeft)
		}
	}

	return result
}
