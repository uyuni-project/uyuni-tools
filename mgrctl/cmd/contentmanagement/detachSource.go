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

type detachSourceFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	ProjectLabel          string
	SourceType            string
	SourceLabel           string
}

func detachSourceCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detachSource",
		Short: "Detach a Source from a Project",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags detachSourceFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, detachSource)
		},
	}

	cmd.Flags().String("ProjectLabel", "", "Content Project label")
	cmd.Flags().String("SourceType", "", "Project Source type, e.g. 'software'")
	cmd.Flags().String("SourceLabel", "", "Project Source label")

	return cmd
}

func detachSource(globalFlags *types.GlobalFlags, flags *detachSourceFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.ProjectLabel, flags.SourceType, flags.SourceLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
