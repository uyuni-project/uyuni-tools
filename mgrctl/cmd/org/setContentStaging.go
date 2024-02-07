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

type setContentStagingFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	OrgId                 int
	Enable                bool
}

func setContentStagingCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setContentStaging",
		Short: "Set the status of content staging for the given organization.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setContentStagingFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setContentStaging)
		},
	}

	cmd.Flags().String("OrgId", "", "")
	cmd.Flags().String("Enable", "", "Use true/false to enable/disable")

	return cmd
}

func setContentStaging(globalFlags *types.GlobalFlags, flags *setContentStagingFlags, cmd *cobra.Command, args []string) error {

	res, err := org.Org(&flags.ConnectionDetails, flags.OrgId, flags.Enable)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
