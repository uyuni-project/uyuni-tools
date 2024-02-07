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

type refreshPillarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sids                  []int
	Subset                string
}

func refreshPillarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refreshPillar",
		Short: "refresh all the pillar data of a list of systems.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags refreshPillarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, refreshPillar)
		},
	}

	cmd.Flags().String("Sids", "", "$desc")
	cmd.Flags().String("Subset", "", "")

	return cmd
}

func refreshPillar(globalFlags *types.GlobalFlags, flags *refreshPillarFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.Sids, flags.Subset)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
