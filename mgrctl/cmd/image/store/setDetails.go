package store

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/image/store"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setDetailsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
}

func setDetailsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setDetails",
		Short: "Set details of an image store",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setDetailsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setDetails)
		},
	}

	cmd.Flags().String("Label", "", "")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setDetails(globalFlags *types.GlobalFlags, flags *setDetailsFlags, cmd *cobra.Command, args []string) error {

res, err := store.Store(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

