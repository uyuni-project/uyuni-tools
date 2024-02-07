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

type buildProjectFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	Message          string
}

func buildProjectCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buildProject",
		Short: "Build a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags buildProjectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, buildProject)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Project label")
	cmd.Flags().String("Message", "", "log message to be assigned to the build")

	return cmd
}

func buildProject(globalFlags *types.GlobalFlags, flags *buildProjectFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.Message)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

