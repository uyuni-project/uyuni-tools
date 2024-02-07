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

type cloneAsyncFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	$param.getFlagName()          $param.getType()
}

func cloneAsyncCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloneAsync",
		Short: "Asynchronously clone a list of errata into the specified channel.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cloneAsyncFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, cloneAsync)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func cloneAsync(globalFlags *types.GlobalFlags, flags *cloneAsyncFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.ChannelLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

