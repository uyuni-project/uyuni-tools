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

type getAllowedOrgsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SlaveId          int
}

func getAllowedOrgsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getAllowedOrgs",
		Short: "Get all orgs this Master is willing to export to the specified Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getAllowedOrgsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getAllowedOrgs)
		},
	}

	cmd.Flags().String("SlaveId", "", "Id of the desired Slave")

	return cmd
}

func getAllowedOrgs(globalFlags *types.GlobalFlags, flags *getAllowedOrgsFlags, cmd *cobra.Command, args []string) error {

res, err := slave.Slave(&flags.ConnectionDetails, flags.SlaveId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

