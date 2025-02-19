// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package backup

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/backup/shared"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func doRestore(
	_ *types.GlobalFlags,
	flags *shared.Flagpole,
	_ *cobra.Command,
	args []string,
) error {
	return nil
}
