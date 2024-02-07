package configchannel

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/configchannel"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateInitSlsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
}

func updateInitSlsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateInitSls",
		Short: "Update the init.sls file for the given state channel. User can only update contents, nothing else.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateInitSlsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateInitSls)
		},
	}

	cmd.Flags().String("Label", "", "the channel label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func updateInitSls(globalFlags *types.GlobalFlags, flags *updateInitSlsFlags, cmd *cobra.Command, args []string) error {

res, err := configchannel.Configchannel(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

