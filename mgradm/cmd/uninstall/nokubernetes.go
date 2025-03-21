// SPDX-FileCopyrightText: 2025 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

//go:build nok8s

package uninstall

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

func uninstallForKubernetes(
	_ *types.GlobalFlags,
	_ *utils.UninstallFlags,
	_ *cobra.Command,
	_ []string,
) error {
	return nil
}
