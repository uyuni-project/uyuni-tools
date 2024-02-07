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

type removeProjectFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
}

func removeProjectCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeProject",
		Short: "Remove Content Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeProjectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeProject)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")

	return cmd
}

func removeProject(globalFlags *types.GlobalFlags, flags *removeProjectFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

