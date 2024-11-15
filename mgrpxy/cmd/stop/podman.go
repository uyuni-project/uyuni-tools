// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package stop

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/podman"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func podmanStop(
	_ *types.GlobalFlags,
	_ *stopFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return podman.StopService(podman.ProxyService)
}
