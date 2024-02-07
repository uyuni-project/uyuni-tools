package slave

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/sync/slave"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setAllowedOrgsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SlaveId          int
	$param.getFlagName()          $param.getType()
}

func setAllowedOrgsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setAllowedOrgs",
		Short: "Set the orgs this Master is willing to export to the specified Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setAllowedOrgsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setAllowedOrgs)
		},
	}

	cmd.Flags().String("SlaveId", "", "ID of the desired Slave")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setAllowedOrgs(globalFlags *types.GlobalFlags, flags *setAllowedOrgsFlags, cmd *cobra.Command, args []string) error {

res, err := slave.Slave(&flags.ConnectionDetails, flags.SlaveId, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

