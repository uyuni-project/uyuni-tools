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

type setMasterOrgsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
	$param.getFlagName()          $param.getType()
}

func setMasterOrgsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setMasterOrgs",
		Short: "Reset all organizations the specified Master has exported to this Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setMasterOrgsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setMasterOrgs)
		},
	}

	cmd.Flags().String("MasterId", "", "Id of the desired Master")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setMasterOrgs(globalFlags *types.GlobalFlags, flags *setMasterOrgsFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

