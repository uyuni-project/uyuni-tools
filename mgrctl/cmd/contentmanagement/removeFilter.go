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

type removeFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	FilterId              int
}

func removeFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeFilter",
		Short: "Remove a Content Filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags removeFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, removeFilter)
		},
	}

	cmd.Flags().String("FilterId", "", "Filter ID")

	return cmd
}

func removeFilter(globalFlags *types.GlobalFlags, flags *removeFilterFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.FilterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
