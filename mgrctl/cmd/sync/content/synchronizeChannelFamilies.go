package content

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/content"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type synchronizeChannelFamiliesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func synchronizeChannelFamiliesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "synchronizeChannelFamilies",
		Short: "Synchronize channel families between the Customer Center
             and the #product() database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags synchronizeChannelFamiliesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, synchronizeChannelFamilies)
		},
	}


	return cmd
}

func synchronizeChannelFamilies(globalFlags *types.GlobalFlags, flags *synchronizeChannelFamiliesFlags, cmd *cobra.Command, args []string) error {

res, err := content.Content(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

