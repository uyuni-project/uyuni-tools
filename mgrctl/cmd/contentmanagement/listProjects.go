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

type listProjectsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listProjectsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProjects",
		Short: "List Content Projects visible to user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProjectsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProjects)
		},
	}

	return cmd
}

func listProjects(globalFlags *types.GlobalFlags, flags *listProjectsFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
