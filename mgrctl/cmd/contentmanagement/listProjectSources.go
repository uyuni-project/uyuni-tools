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

type listProjectSourcesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
}

func listProjectSourcesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProjectSources",
		Short: "List Content Project Sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProjectSourcesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProjectSources)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")

	return cmd
}

func listProjectSources(globalFlags *types.GlobalFlags, flags *listProjectSourcesFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
