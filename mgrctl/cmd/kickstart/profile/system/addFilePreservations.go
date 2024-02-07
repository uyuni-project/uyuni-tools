package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/profile/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addFilePreservationsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	$param.getFlagName()          $param.getType()
}

func addFilePreservationsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addFilePreservations",
		Short: "Adds the given list of file preservations to the specified kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addFilePreservationsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addFilePreservations)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func addFilePreservations(globalFlags *types.GlobalFlags, flags *addFilePreservationsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

