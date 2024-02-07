package snippet

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/kickstart/snippet"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type listAllFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listAllCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listAll",
		Short: "List all cobbler snippets for the logged in user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listAllFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listAll)
		},
	}

	return cmd
}

func listAll(globalFlags *types.GlobalFlags, flags *listAllFlags, cmd *cobra.Command, args []string) error {

	res, err := snippet.Snippet(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}
