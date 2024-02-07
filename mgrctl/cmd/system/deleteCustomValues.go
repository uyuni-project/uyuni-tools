package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type deleteCustomValuesFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Keys          []string
}

func deleteCustomValuesCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteCustomValues",
		Short: "Delete the custom values defined for the custom system information keys
 provided from the given system.
 (Note: Attempt to delete values of non-existing keys throws exception. Attempt to
 delete value of existing key which has assigned no values doesn't throw exception.)",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteCustomValuesFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteCustomValues)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Keys", "", "$desc")

	return cmd
}

func deleteCustomValues(globalFlags *types.GlobalFlags, flags *deleteCustomValuesFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Keys)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

