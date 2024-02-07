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

type updateEnvironmentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	EnvLabel          string
	$param.getFlagName()          $param.getType()
}

func updateEnvironmentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateEnvironment",
		Short: "Update Content Environment with given label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateEnvironmentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateEnvironment)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("EnvLabel", "", "Content Environment label")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func updateEnvironment(globalFlags *types.GlobalFlags, flags *updateEnvironmentFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.EnvLabel, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

