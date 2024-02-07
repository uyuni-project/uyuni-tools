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

type getClmSyncPatchesConfigFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func getClmSyncPatchesConfigCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getClmSyncPatchesConfig",
		Short: "Reads the content lifecycle management patch synchronization config option.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getClmSyncPatchesConfigFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getClmSyncPatchesConfig)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func getClmSyncPatchesConfig(globalFlags *types.GlobalFlags, flags *getClmSyncPatchesConfigFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

