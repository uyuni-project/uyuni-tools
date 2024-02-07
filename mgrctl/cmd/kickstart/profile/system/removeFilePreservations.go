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

type removeFilePreservationsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
	$param.getFlagName()          $param.getType()
}

func removeFilePreservationsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeFilePreservations",
		Short: "Removes the given list of file preservations from the specified
 kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeFilePreservationsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeFilePreservations)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func removeFilePreservations(globalFlags *types.GlobalFlags, flags *removeFilePreservationsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

