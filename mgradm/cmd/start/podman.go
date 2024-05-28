// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package start

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanStart(
	globalFlags *types.GlobalFlags,
	flags *startFlags,
	cmd *cobra.Command,
	args []string,
) error {
	err1 := podman.StartInstantiated(podman.ServerAttestationService)
	err2 := podman.StartInstantiated(podman.HubXmlrpcService)
	err3 := podman.StartService(podman.ServerService)
	return errors.Join(err1, err2, err3)
}
