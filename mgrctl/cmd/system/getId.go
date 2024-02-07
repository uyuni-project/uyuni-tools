package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type getIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Name          string
}

func getIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getId",
		Short: "Get system IDs and last check in information for the given system name.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getId)
		},
	}

	cmd.Flags().String("Name", "", "")

	return cmd
}

func getId(globalFlags *types.GlobalFlags, flags *getIdFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

