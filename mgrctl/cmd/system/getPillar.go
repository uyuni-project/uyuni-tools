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

type getPillarFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SystemId              int
	Category              string
	MinionId              int
}

func getPillarCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getPillar",
		Short: "Get pillar data of given category for given system",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags getPillarFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, getPillar)
		},
	}

	cmd.Flags().String("SystemId", "", "")
	cmd.Flags().String("Category", "", "")
	cmd.Flags().String("MinionId", "", "")

	return cmd
}

func getPillar(globalFlags *types.GlobalFlags, flags *getPillarFlags, cmd *cobra.Command, args []string) error {

	res, err := system.System(&flags.ConnectionDetails, flags.SystemId, flags.Category, flags.MinionId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
