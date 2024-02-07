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

type listCustomFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
}

func listCustomCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listCustom",
		Short: "List only custom snippets for the logged in user.
    These snipppets are editable.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags listCustomFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, listCustom)
		},
	}


	return cmd
}

func listCustom(globalFlags *types.GlobalFlags, flags *listCustomFlags, cmd *cobra.Command, args []string) error {

res, err := snippet.Snippet(&flags.ConnectionDetails)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

