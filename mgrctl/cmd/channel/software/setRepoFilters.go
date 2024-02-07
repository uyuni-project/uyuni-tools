package software

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/channel/software"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type setRepoFiltersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
	$param.getFlagName()          $param.getType()
}

func setRepoFiltersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setRepoFilters",
		Short: "Replaces the existing set of filters for a given repo.
 Filters are ranked by their order in the array.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags setRepoFiltersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, setRepoFilters)
		},
	}

	cmd.Flags().String("Label", "", "repository label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func setRepoFilters(globalFlags *types.GlobalFlags, flags *setRepoFiltersFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Label, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

