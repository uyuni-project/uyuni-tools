// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/types"
)

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "org",
		Short: "Organization-related commands",
	}

	api.AddAPIFlags(cmd, false)

	cmd.AddCommand(getClmSyncPatchesConfigCommand(globalFlags))
	cmd.AddCommand(transferSystemsCommand(globalFlags))
	cmd.AddCommand(setErrataEmailNotifsForOrgCommand(globalFlags))
	cmd.AddCommand(setContentStagingCommand(globalFlags))
	cmd.AddCommand(isContentStagingEnabledCommand(globalFlags))
	cmd.AddCommand(getDetailsCommand(globalFlags))
	cmd.AddCommand(updateNameCommand(globalFlags))
	cmd.AddCommand(setPolicyForScapFileUploadCommand(globalFlags))
	cmd.AddCommand(isErrataEmailNotifsForOrgCommand(globalFlags))
	cmd.AddCommand(deleteCommand(globalFlags))
	cmd.AddCommand(setOrgConfigManagedByOrgAdminCommand(globalFlags))
	cmd.AddCommand(getPolicyForScapFileUploadCommand(globalFlags))
	cmd.AddCommand(listUsersCommand(globalFlags))
	cmd.AddCommand(setPolicyForScapResultDeletionCommand(globalFlags))
	cmd.AddCommand(setClmSyncPatchesConfigCommand(globalFlags))
	cmd.AddCommand(isOrgConfigManagedByOrgAdminCommand(globalFlags))
	cmd.AddCommand(createCommand(globalFlags))
	cmd.AddCommand(getPolicyForScapResultDeletionCommand(globalFlags))
	cmd.AddCommand(createFirstCommand(globalFlags))
	cmd.AddCommand(listOrgsCommand(globalFlags))

	return cmd
}
