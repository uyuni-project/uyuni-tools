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

type lookupFilterFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	FilterId              int
}

func lookupFilterCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupFilter",
		Short: "Lookup a Content Filter by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupFilterFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupFilter)
		},
	}

	cmd.Flags().String("FilterId", "", "Filter ID")

	return cmd
}

func lookupFilter(globalFlags *types.GlobalFlags, flags *lookupFilterFlags, cmd *cobra.Command, args []string) error {

	res, err := contentmanagement.Contentmanagement(&flags.ConnectionDetails, flags.FilterId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
