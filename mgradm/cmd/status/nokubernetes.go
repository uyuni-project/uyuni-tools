// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package status

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kubernetesStatus(
	globalFlags *types.GlobalFlags,
	flags *statusFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return fmt.Errorf("built without kubernetes support")
}
