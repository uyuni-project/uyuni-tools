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

type getNetworkForSystemsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids                  []int
}

func getNetworkForSystemsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getNetworkForSystems",
		Short: "Get the addresses and hostname for a given list of systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getNetworkForSystemsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getNetworkForSystems)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")

	return cmd
}

func getNetworkForSystems(globalFlags *types.GlobalFlags, flags *getNetworkForSystemsFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sids)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
