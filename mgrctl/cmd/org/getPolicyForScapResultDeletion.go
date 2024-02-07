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

type getPolicyForScapResultDeletionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func getPolicyForScapResultDeletionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPolicyForScapResultDeletion",
		Short: "Get the status of SCAP result deletion settings for the given
 organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPolicyForScapResultDeletionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPolicyForScapResultDeletion)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func getPolicyForScapResultDeletion(globalFlags *types.GlobalFlags, flags *getPolicyForScapResultDeletionFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

