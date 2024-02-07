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

type createProjectFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	Name                  string
	Description           string
}

func createProjectCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createProject",
		Short: "Create Content Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createProjectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createProject)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("Name", "", "Content Project name")
	cmd.Flags().String("Description", "", "Content Project description")

	return cmd
}

func createProject(globalFlags *types.GlobalFlags, flags *createProjectFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.Name, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
