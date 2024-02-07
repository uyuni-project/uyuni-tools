package highstate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/recurring/highstate"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type updateFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	$param.getFlagName()          $param.getType()
}

func updateCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the properties of a recurring highstate action.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, update)
		},
	}

	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func update(globalFlags *types.GlobalFlags, flags *updateFlags, cmd *cobra.Command, args []string) error {

res, err := highstate.Highstate(&flags.ConnectionDetails, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

