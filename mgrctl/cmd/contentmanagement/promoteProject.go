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

type promoteProjectFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	EnvLabel          string
}

func promoteProjectCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promoteProject",
		Short: "Promote an Environment in a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags promoteProjectFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, promoteProject)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Project label")
	cmd.Flags().String("EnvLabel", "", "Environment label")

	return cmd
}

func promoteProject(globalFlags *types.GlobalFlags, flags *promoteProjectFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.EnvLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

