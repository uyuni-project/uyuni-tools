// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package start

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kubernetesStart(
	globalFlags *types.GlobalFlags,
	flags *startFlags,
	cmd *cobra.Command,
	args []string,
) error {
	return fmt.Errorf("built without kubernetes support")
}
