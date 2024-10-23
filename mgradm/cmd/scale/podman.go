// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package scale

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"

	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
)

var systemd podman.Systemd = podman.SystemdImpl{}

func podmanScale(
	globalFlags *types.GlobalFlags,
	flags *scaleFlags,
	cmd *cobra.Command,
	args []string,
) error {
	newReplicas := flags.Replicas
	service := args[0]
	if service == podman.ServerAttestationService {
		return systemd.ScaleService(newReplicas, service)
	}
	if service == podman.HubXmlrpcService {
		if newReplicas > 1 {
			return errors.New(L("Multiple Hub XML-RPC container replicas are not currently supported."))
		}
		return systemd.ScaleService(newReplicas, service)
	}
	return fmt.Errorf(L("service not allowing to be scaled: %s"), service)
}
