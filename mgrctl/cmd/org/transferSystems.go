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

type transferSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ToOrgId          int
	Sids          []int
}

func transferSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transferSystems",
		Short: "Transfer systems from one organization to another.  If executed by
 a #product() administrator, the systems will be transferred from their current
 organization to the organization specified by the toOrgId.  If executed by
 an organization administrator, the systems must exist in the same organization
 as that administrator and the systems will be transferred to the organization
 specified by the toOrgId. In any scenario, the origination and destination
 organizations must be defined in a trust.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags transferSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, transferSystems)
		},
	}

	cmd.Flags().String("ToOrgId", "", "ID of the organization where the system(s) will be transferred to.")
	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func transferSystems(globalFlags *types.GlobalFlags, flags *transferSystemsFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.ToOrgId, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

