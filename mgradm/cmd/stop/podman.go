// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanStop(
	globalFlags *types.GlobalFlags,
	flags *stopFlags,
	cmd *cobra.Command,
	args []string,
) error {
	err1 := podman.StopInstantiated(podman.ServerAttestationService)
	err2 := podman.StopInstantiated(podman.HubXmlrpcService)
	err3 := podman.StopService(podman.ServerService)
	return errors.Join(err1, err2, err3)
}
