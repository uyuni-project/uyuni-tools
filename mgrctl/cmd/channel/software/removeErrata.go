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

type removeErrataFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	$param.getFlagName()          $param.getType()
	RemovePackages          bool
}

func removeErrataCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeErrata",
		Short: "Removes a given list of errata from the given channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeErrataFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeErrata)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "target channel")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")
	cmd.Flags().String("RemovePackages", "", "true to remove packages from the channel")

	return cmd
}

func removeErrata(globalFlags *types.GlobalFlags, flags *removeErrataFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.ChannelLabel, flags.$param.getFlagName(), flags.RemovePackages)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

