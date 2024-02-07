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

type createAppStreamFiltersFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Prefix                string
	ChannelLabel          string
	ProjectLabel          string
}

func createAppStreamFiltersCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createAppStreamFilters",
		Short: "Create Filters for AppStream Modular Channel and attach them to CLM Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createAppStreamFiltersFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createAppStreamFilters)
		},
	}

	cmd.Flags().String("Prefix", "", "Filter name prefix")
	cmd.Flags().String("ChannelLabel", "", "Modular Channel label")
	cmd.Flags().String("ProjectLabel", "", "Project label")

	return cmd
}

func createAppStreamFilters(globalFlags *types.GlobalFlags, flags *createAppStreamFiltersFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.Prefix, flags.ChannelLabel, flags.ProjectLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
