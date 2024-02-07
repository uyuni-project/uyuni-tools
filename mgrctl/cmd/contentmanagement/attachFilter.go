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

type attachFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	FilterId              int
}

func attachFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachFilter",
		Short: "Attach a Filter to a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags attachFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, attachFilter)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Project label")
	cmd.Flags().String("FilterId", "", "filter ID to attach")

	return cmd
}

func attachFilter(globalFlags *types.GlobalFlags, flags *attachFilterFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.FilterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
