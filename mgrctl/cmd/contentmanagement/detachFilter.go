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

type detachFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	FilterId              int
}

func detachFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detachFilter",
		Short: "Detach a Filter from a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags detachFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, detachFilter)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Project label")
	cmd.Flags().String("FilterId", "", "filter ID to detach")

	return cmd
}

func detachFilter(globalFlags *types.GlobalFlags, flags *detachFilterFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.FilterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
