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

type isContentStagingEnabledFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId          int
}

func isContentStagingEnabledCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isContentStagingEnabled",
		Short: "Get the status of content staging settings for the given organization.
 Returns true if enabled, false otherwise.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags isContentStagingEnabledFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, isContentStagingEnabled)
		},
	}

	cmd.Flags().String("OrgId", "", "")

	return cmd
}

func isContentStagingEnabled(globalFlags *types.GlobalFlags, flags *isContentStagingEnabledFlags, cmd *cobra.Command, args []string) error {

res, err := org.Org(&flags.ConnectionDetails, flags.OrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

