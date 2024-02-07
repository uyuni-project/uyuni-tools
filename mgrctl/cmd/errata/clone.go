package errata

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/errata"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type cloneFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	$param.getFlagName()          $param.getType()
}

func cloneCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Clone a list of errata into the specified channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cloneFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, clone)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func clone(globalFlags *types.GlobalFlags, flags *cloneFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.ChannelLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

