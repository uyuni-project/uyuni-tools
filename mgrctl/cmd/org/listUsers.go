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

type listUsersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func listUsersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listUsers",
		Short: "Returns the list of users in a given organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listUsersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listUsers)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func listUsers(globalFlags *types.GlobalFlags, flags *listUsersFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

