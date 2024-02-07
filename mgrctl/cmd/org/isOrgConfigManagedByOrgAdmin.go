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

type isOrgConfigManagedByOrgAdminFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func isOrgConfigManagedByOrgAdminCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isOrgConfigManagedByOrgAdmin",
		Short: "Returns whether Organization Administrator is able to manage his
 organization configuration. This may have a high impact on general #product() performance.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isOrgConfigManagedByOrgAdminFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isOrgConfigManagedByOrgAdmin)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func isOrgConfigManagedByOrgAdmin(globalFlags *types.GlobalFlags, flags *isOrgConfigManagedByOrgAdminFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

