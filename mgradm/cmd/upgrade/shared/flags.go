// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	"github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	shared_utils "github.com/uyuni-project/uyuni-tools/shared/utils"
)

// UpgradeFlags represents flags used for upgrading a server.
type UpgradeFlags struct {
	Image          types.ImageFlags `mapstructure:",squash"`
	DbUpgradeImage types.ImageFlags `mapstructure:"dbupgrade"`
	Coco           shared.CocoFlags
}

// AddUpgradeFlags add upgrade flags to a command.
func AddUpgradeFlags(cmd *cobra.Command) {
	utils.AddImageFlag(cmd)
	utils.AddDbUpgradeImageFlag(cmd)

	_ = shared_utils.AddFlagHelpGroup(cmd, &shared_utils.Group{
		ID:    "coco-container",
		Title: L("Confidential Computing Flags"),
	})
	utils.AddContainerImageFlags(cmd, "coco", L("confidential computing attestation"), "coco-container", "server-attestation")
}

// AddUpgradeListFlags add upgrade list flags to a command.
func AddUpgradeListFlags(cmd *cobra.Command) {
	utils.AddImageFlag(cmd)
}
