package recurring

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/recurring"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type lookupByIdFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Id          int
}

func lookupByIdCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lookupById",
		Short: "Find a recurring action with the given action ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags lookupByIdFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, lookupById)
		},
	}

	cmd.Flags().String("Id", "", "the action ID")

	return cmd
}

func lookupById(globalFlags *types.GlobalFlags, flags *lookupByIdFlags, cmd *cobra.Command, args []string) error {

res, err := recurring.Recurring(&flags.ConnectionDetails, flags.Id)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

