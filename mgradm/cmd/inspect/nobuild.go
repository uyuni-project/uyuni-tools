// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package inspect

import (
	"github.com/spf13/cobra"
	inspect_shared "github.com/uyuni-project/uyuni-tools/mgradm/cmd/inspect/shared"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kuberneteInspect(
	globalFlags *types.GlobalFlags,
	flags *inspect_shared.InspectFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return nil
}
