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

type mapToLocalFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	MasterId          int
	MasterOrgId          int
	LocalOrgId          int
}

func mapToLocalCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mapToLocal",
		Short: "Add a single organizations to the list of those the specified Master has
 exported to this Slave",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags mapToLocalFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, mapToLocal)
		},
	}

	cmd.Flags().String("MasterId", "", "ID of the desired Master")
	cmd.Flags().String("MasterOrgId", "", "ID of the desired Master")
	cmd.Flags().String("LocalOrgId", "", "ID of the desired Master")

	return cmd
}

func mapToLocal(globalFlags *types.GlobalFlags, flags *mapToLocalFlags, cmd *cobra.Command, args []string) error {

res, err := master.Master(&flags.ConnectionDetails, flags.MasterId, flags.MasterOrgId, flags.LocalOrgId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

