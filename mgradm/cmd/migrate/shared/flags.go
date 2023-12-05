// SPDX-FileCopyrightText: 2023 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
)

type MigrateFlags struct {
	Image cmd_utils.ImageFlags `mapstructure:",squash"`
}

func AddMigrateFlags(cmd *cobra.Command) {
	cmd_utils.AddImageFlag(cmd)
}
