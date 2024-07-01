// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package shared

import (
	"path"

	"github.com/spf13/cobra"
	cmd_utils "github.com/uyuni-project/uyuni-tools/mgradm/shared/utils"
	. "github.com/uyuni-project/uyuni-tools/shared/l10n"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

// MigrateFlags represents flag required by migration command.
type MigrateFlags struct {
	Image          types.ImageFlags `mapstructure:",squash"`
	DbUpgradeImage types.ImageFlags `mapstructure:"dbupgrade"`
	HubXmlrpc      types.HubXmlrpcFlags
	User           string
	Mirror         string
}

// AddMigrateFlags add migration flags to a command.
func AddMigrateFlags(cmd *cobra.Command) {
	cmd_utils.AddMirrorFlag(cmd)
	cmd_utils.AddImageFlag(cmd)
	cmd_utils.AddDbUpgradeImageFlag(cmd)
	cmd.Flags().String("user", "root", L("User on the source server. Non-root user must have passwordless sudo privileges (NOPASSWD tag in /etc/sudoers)."))

	cmd.Flags().Int("hubxmlrpc-replicas", 0, L("How many replicas of the Hub XML-RPC API service container should be started. (only 0 or 1 supported for now)"))
	hubXmlrpcImage := path.Join(utils.DefaultNamespace, "server-hub-xmlrpc-api")
	cmd.Flags().String("hubxmlrpc-image", hubXmlrpcImage, L("Hub XML-RPC API Image"))
	cmd.Flags().String("hubxmlrpc-tag", utils.DefaultTag, L("Hub XML-RPC API Image Tag"))

	_ = utils.AddFlagHelpGroup(cmd, &utils.Group{ID: "hubxmlrpc-container", Title: L("Hub XML-RPC API")})
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-replicas", "hubxmlrpc-container")
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-image", "hubxmlrpc-container")
	_ = utils.AddFlagToHelpGroupID(cmd, "hubxmlrpc-tag", "hubxmlrpc-container")
}
