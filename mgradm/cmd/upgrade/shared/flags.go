// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"path"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/mgradm/cmd/install/shared"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// UpgradeFlags represents flags used for upgrading a server.
type UpgradeFlags struct {
	Image          types.ImageFlags `mapstructure:",squash"`
	DbUpgradeImage types.ImageFlags `mapstructure:"dbupgrade"`
	HubXmlrpc      types.HubXmlrpcFlags
	Coco           shared.CocoFlags
}

// AddUpgradeFlags add upgrade flags to a command.
func AddUpgradeFlags(cmd *cobra.Command) {
	cmd_utils.AddImageFlag(cmd)
	cmd_utils.AddDbUpgradeImageFlag(cmd)
	cmd_utils.AddContainerImageFlags(cmd, "coco", L("confidential computing attestation"))

	hubXmlrpcImage := path.Join(utils.DefaultNamespace, "server-hub-xmlrpc-api")
	cmd.Flags().String("hubxmlrpc-image", hubXmlrpcImage, L("Hub XML-RPC API Image"))
	cmd.Flags().String("hubxmlrpc-tag", utils.DefaultTag, L("Hub XML-RPC API Image Tag"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-image", "hubxmlrpc-container")
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-tag", "hubxmlrpc-container")
}

// AddUpgradeListFlags add upgrade list flags to a command.
func AddUpgradeListFlags(cmd *cobra.Command) {
	cmd_utils.AddImageFlag(cmd)
}
