package org

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/org"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type disableAccessFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	OrgId          int
}

func disableAccessCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disableAccess",
		Short: "Disable access to the channel for the given organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags disableAccessFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, disableAccess)
		},
	}

	cmd.Flags().String("Label", "", "label of the channel")
	cmd.Flags().String("OrgId", "", "ID of org being removed access")

	return cmd
}

func disableAccess(globalFlags *types.GlobalFlags, flags *disableAccessFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.Label, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

