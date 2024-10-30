// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package inspect

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func kuberneteInspect(
	_ *types.GlobalFlags,
	_ *inspectFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return nil
}
