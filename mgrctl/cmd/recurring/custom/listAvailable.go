package custom

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/recurring/custom"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAvailableFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAvailableCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAvailable",
		Short: "List all the custom states available to the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAvailableFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAvailable)
		},
	}

	return cmd
}

func listAvailable(globalFlags *types.GlobalFlags, flags *listAvailableFlags, cmd *cobra.Command, args []string) error {

	res, err := custom.Custom(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
