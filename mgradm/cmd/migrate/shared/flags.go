// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

// MigrateFlags represents flag required by migration command.
type MigrateFlags struct {
	Image          types.ImageFlags `mapstructure:",squash"`
	MigrationImage types.ImageFlags `mapstructure:"migration"`
	User           string
	Mirror         string
}

// AddMigrateFlags add migration flags to a command.
func AddMigrateFlags(cmd *cobra.Command) {
	utils.AddMirrorFlag(cmd)
	utils.AddImageFlag(cmd)
	utils.AddMigrationImageFlag(cmd)
	cmd.Flags().String("user", "root", L("User on the source server. Non-root user must have passwordless sudo privileges (NOPASSWD tag in /etc/sudoers)."))
}
