package contentmanagement

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/contentmanagement"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listFiltersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listFiltersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listFilters",
		Short: "List all Content Filters visible to given user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listFiltersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listFilters)
		},
	}


	return cmd
}

func listFilters(globalFlags *types.GlobalFlags, flags *listFiltersFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

