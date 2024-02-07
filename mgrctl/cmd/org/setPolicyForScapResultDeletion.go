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

type setPolicyForScapResultDeletionFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func setPolicyForScapResultDeletionCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setPolicyForScapResultDeletion",
		Short: "Set the status of SCAP result deletion settins for the given
 organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setPolicyForScapResultDeletionFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setPolicyForScapResultDeletion)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func setPolicyForScapResultDeletion(globalFlags *types.GlobalFlags, flags *setPolicyForScapResultDeletionFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

