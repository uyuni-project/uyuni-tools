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

type updateFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	FilterId          int
	Name          string
	Rule          string
	$param.getFlagName()          $param.getType()
}

func updateFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "updateFilter",
		Short: "Update a Content Filter
 #paragraph_end()
 #paragraph()
 See also: createFilter(), listFilterCriteria()",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags updateFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, updateFilter)
		},
	}

	cmd.Flags().String("FilterId", "", "Filter ID")
	cmd.Flags().String("Name", "", "New filter name")
	cmd.Flags().String("Rule", "", "New filter rule ('deny' or 'allow')")
	cmd.Flags().String("$param.getFlagName()", "", "$param.getDesc()")

	return cmd
}

func updateFilter(globalFlags *types.GlobalFlags, flags *updateFilterFlags, cmd *cobra.Command, args []string) error {

res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.FilterId, flags.Name, flags.Rule, flags.$param.getFlagName())
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

