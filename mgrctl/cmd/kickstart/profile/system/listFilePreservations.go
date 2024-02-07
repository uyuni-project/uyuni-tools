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

type listFilePreservationsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	KsLabel          string
}

func listFilePreservationsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFilePreservations",
		Short: "Returns the set of all file preservations associated with the given
 kickstart profile.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFilePreservationsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFilePreservations)
		},
	}

	cmd.Flags().String("KsLabel", "", "the kickstart profile label")

	return cmd
}

func listFilePreservations(globalFlags *types.GlobalFlags, flags *listFilePreservationsFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.KsLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

