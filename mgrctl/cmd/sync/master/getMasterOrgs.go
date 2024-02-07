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

type getMasterOrgsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
}

func getMasterOrgsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getMasterOrgs",
		Short: "List all organizations the specified Master has exported to this Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getMasterOrgsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getMasterOrgs)
		},
	}

	cmd.Flags().String("MasterId", "", "ID of the desired Master")

	return cmd
}

func getMasterOrgs(globalFlags *types.GlobalFlags, flags *getMasterOrgsFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

