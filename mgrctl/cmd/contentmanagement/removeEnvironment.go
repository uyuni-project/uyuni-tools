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

type removeEnvironmentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	EnvLabel              string
}

func removeEnvironmentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeEnvironment",
		Short: "Remove a Content Environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeEnvironmentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeEnvironment)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("EnvLabel", "", "Content Environment label")

	return cmd
}

func removeEnvironment(globalFlags *types.GlobalFlags, flags *removeEnvironmentFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.EnvLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
