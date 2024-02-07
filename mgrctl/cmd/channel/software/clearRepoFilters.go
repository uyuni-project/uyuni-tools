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

type clearRepoFiltersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Label          string
}

func clearRepoFiltersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clearRepoFilters",
		Short: "Removes the filters for a repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags clearRepoFiltersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, clearRepoFilters)
		},
	}

	cmd.Flags().String("Label", "", "repository label")

	return cmd
}

func clearRepoFilters(globalFlags *types.GlobalFlags, flags *clearRepoFiltersFlags, cmd *cobra.Command, args []string) error {

res, err := software.Software(&flags.ConnectionDetails, flags.Label)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

