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

type listProjectFiltersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
}

func listProjectFiltersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProjectFilters",
		Short: "List all Filters associated with a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProjectFiltersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProjectFilters)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Project label")

	return cmd
}

func listProjectFilters(globalFlags *types.GlobalFlags, flags *listProjectFiltersFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

