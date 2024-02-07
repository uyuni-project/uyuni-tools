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

type cloneAsOriginalAsyncFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ChannelLabel          string
	$param.getFlagName()          $param.getType()
}

func cloneAsOriginalAsyncCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloneAsOriginalAsync",
		Short: "Asynchronously clones a list of errata into a specified cloned channel
 according the original erratas",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags cloneAsOriginalAsyncFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, cloneAsOriginalAsync)
		},
	}

	cmd.Flags().String("ChannelLabel", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func cloneAsOriginalAsync(globalFlags *types.GlobalFlags, flags *cloneAsOriginalAsyncFlags, cmd *cobra.Command, args []string) error {

res, err := errata.Errata(&flags.ConnectionDetails, flags.ChannelLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

