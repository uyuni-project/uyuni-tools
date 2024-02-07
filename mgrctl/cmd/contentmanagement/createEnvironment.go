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

type createEnvironmentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	PredecessorLabel          string
	EnvLabel          string
	Name          string
	Description          string
}

func createEnvironmentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createEnvironment",
		Short: "Create a Content Environment and appends it behind given Content Environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags createEnvironmentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, createEnvironment)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("PredecessorLabel", "", "Predecessor Environment label")
	cmd.Flags().String("EnvLabel", "", "new Content Environment label")
	cmd.Flags().String("Name", "", "new Content Environment name")
	cmd.Flags().String("Description", "", "new Content Environment description")

	return cmd
}

func createEnvironment(globalFlags *types.GlobalFlags, flags *createEnvironmentFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.PredecessorLabel, flags.EnvLabel, flags.Name, flags.Description)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

