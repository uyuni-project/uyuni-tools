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

type listProjectEnvironmentsFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
}

func listProjectEnvironmentsCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listProjectEnvironments",
		Short: "List Environments in a Content Project with the respect to their ordering",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listProjectEnvironmentsFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listProjectEnvironments)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")

	return cmd
}

func listProjectEnvironments(globalFlags *types.GlobalFlags, flags *listProjectEnvironmentsFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

