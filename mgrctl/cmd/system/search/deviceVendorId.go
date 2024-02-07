package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system/search"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deviceVendorIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	SearchTerm            string
}

func deviceVendorIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deviceVendorId",
		Short: "List the systems which match this device vendor_id",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deviceVendorIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deviceVendorId)
		},
	}

	cmd.Flags().String("SearchTerm", "", "")

	return cmd
}

func deviceVendorId(globalFlags *types.GlobalFlags, flags *deviceVendorIdFlags, cmd *cobra.Command, args []string) error {

	res, err := search.Search(&flags.ConnectionDetails, flags.SearchTerm)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
