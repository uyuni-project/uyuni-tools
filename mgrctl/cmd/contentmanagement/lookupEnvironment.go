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

type lookupEnvironmentFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	EnvLabel              string
}

func lookupEnvironmentCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupEnvironment",
		Short: "Look up Content Environment based on Content Project and Content Environment label",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupEnvironmentFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupEnvironment)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("EnvLabel", "", "Content Environment label")

	return cmd
}

func lookupEnvironment(globalFlags *types.GlobalFlags, flags *lookupEnvironmentFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.EnvLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
