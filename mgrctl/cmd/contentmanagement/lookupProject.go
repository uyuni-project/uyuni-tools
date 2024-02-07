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

type lookupProjectFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
}

func lookupProjectCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupProject",
		Short: "Look up Content Project with given label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupProjectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupProject)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")

	return cmd
}

func lookupProject(globalFlags *types.GlobalFlags, flags *lookupProjectFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
