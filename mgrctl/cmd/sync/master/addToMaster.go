package master

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/master"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addToMasterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
	$param.getFlagName()          $param.getType()
}

func addToMasterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addToMaster",
		Short: "Add a single organizations to the list of those the specified Master has
 exported to this Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addToMasterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addToMaster)
		},
	}

	cmd.Flags().String("MasterId", "", "Id of the desired Master")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func addToMaster(globalFlags *types.GlobalFlags, flags *addToMasterFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

