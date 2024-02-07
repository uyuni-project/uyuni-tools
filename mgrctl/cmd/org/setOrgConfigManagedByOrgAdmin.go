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

type setOrgConfigManagedByOrgAdminFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
	Enable          bool
}

func setOrgConfigManagedByOrgAdminCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setOrgConfigManagedByOrgAdmin",
		Short: "Sets whether Organization Administrator can manage his organization
 configuration. This may have a high impact on general #product() performance.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setOrgConfigManagedByOrgAdminFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setOrgConfigManagedByOrgAdmin)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("Enable", "", "Use true/false to enable/disable")

	return cmd
}

func setOrgConfigManagedByOrgAdmin(globalFlags *types.GlobalFlags, flags *setOrgConfigManagedByOrgAdminFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId, flags.Enable)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

