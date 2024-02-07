package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type removeRepoFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
}

func removeRepoFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeRepoFilter",
		Short: "Removes a filter for a given repo.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeRepoFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeRepoFilter)
		},
	}

	cmd.Flags().String("Label", "", "repository label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func removeRepoFilter(globalFlags *types.GlobalFlags, flags *removeRepoFilterFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

