// SPDX-FileCopyrightText: 2024 SUSE LLC
//
// SPDX-License-Identifier: Apache-2.0

package org

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/org"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type createFirstFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Organization          string
	Admin                 apiTypes.User
}

func createFirstCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createFirst",
		Short: "create the first user and organization",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createFirstFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createFirst)
		},
	}

	cmd.Flags().String("admin-login", "admin", "Administrator user name")
	cmd.Flags().String("admin-password", "", "Administrator password. If empty, the first user will not be created")
	cmd.Flags().String("admin-firstName", "Administrator", "The first name of the administrator")
	cmd.Flags().String("admin-lastName", "McAdmin", "The last name of the administrator")
	cmd.Flags().String("admin-email", "root@localhost", "The administrator's email")
	cmd.Flags().String("organization", "Organiszation", "The first organization name")

	return cmd
}

func createFirst(globalFlags *types.GlobalFlags, flags *createFirstFlags, cmd *cobra.Command, args []string) error {
	org, err := org.CreateFirst(&flags.ConnectionDetails, flags.Organization, &flags.Admin)
	if err != nil {
		return err
	}

	fmt.Printf("Organization %s created with id %d", org.Name, org.Id)

	return nil
}
