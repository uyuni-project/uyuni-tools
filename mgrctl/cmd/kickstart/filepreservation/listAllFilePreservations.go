package filepreservation

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/filepreservation"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAllFilePreservationsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllFilePreservationsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAllFilePreservations",
		Short: "List all file preservation lists for the organization
 associated with the user logged into the given session",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllFilePreservationsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAllFilePreservations)
		},
	}


	return cmd
}

func listAllFilePreservations(globalFlags *types.GlobalFlags, flags *listAllFilePreservationsFlags, cmd *cobra.Command, args []string) error {

res, err := filepreservation.Filepreservation(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

