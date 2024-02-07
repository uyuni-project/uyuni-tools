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

type attachSourceFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	SourceType          string
	SourceLabel          string
	SourcePosition          int
}

func attachSourceCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachSource",
		Short: "Attach a Source to a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags attachSourceFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, attachSource)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("SourceType", "", "Project Source type, e.g. 'software'")
	cmd.Flags().String("SourceLabel", "", "Project Source label")
	cmd.Flags().String("SourcePosition", "", "Project Source position")

	return cmd
}

func attachSource(globalFlags *types.GlobalFlags, flags *attachSourceFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.SourceType, flags.SourceLabel, flags.SourcePosition)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

