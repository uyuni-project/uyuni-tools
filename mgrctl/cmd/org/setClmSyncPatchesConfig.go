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

type setClmSyncPatchesConfigFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
	Value          bool
}

func setClmSyncPatchesConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setClmSyncPatchesConfig",
		Short: "Sets the content lifecycle management patch synchronization config option.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setClmSyncPatchesConfigFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setClmSyncPatchesConfig)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("Value", "", "The config option value")

	return cmd
}

func setClmSyncPatchesConfig(globalFlags *types.GlobalFlags, flags *setClmSyncPatchesConfigFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId, flags.Value)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

